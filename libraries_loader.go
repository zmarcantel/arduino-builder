/*
 * This file is part of Arduino Builder.
 *
 * Arduino Builder is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA
 *
 * As a special exception, you may use this file as part of a free software
 * library without restriction.  Specifically, if other files instantiate
 * templates or use macros or inline functions from this file, or you compile
 * this file and link it with other files to produce an executable, this
 * file does not by itself cause the resulting executable to be covered by
 * the GNU General Public License.  This exception does not however
 * invalidate any other reasons why the executable file might be covered by
 * the GNU General Public License.
 *
 * Copyright 2015 Arduino LLC (http://www.arduino.cc/)
 */

package builder

import (
	"os"

	"github.com/arduino/arduino-builder/i18n"
	"github.com/arduino/arduino-builder/types"
	"github.com/arduino/arduino-builder/utils"
	"github.com/arduino/go-paths-helper"
	"github.com/bcmi-labs/arduino-cli/arduino/libraries"
)

type LibrariesLoader struct{}

func (s *LibrariesLoader) Run(ctx *types.Context) error {
	builtInLibrariesFolders := ctx.BuiltInLibrariesFolders
	if err := builtInLibrariesFolders.ToAbs(); err != nil {
		return i18n.WrapError(err)
	}
	sortedLibrariesFolders := builtInLibrariesFolders.Clone()

	platform := ctx.TargetPlatform
	debugLevel := ctx.DebugLevel
	logger := ctx.GetLogger()

	actualPlatform := ctx.ActualPlatform
	if actualPlatform != platform {
		actualPlatformLibDir := paths.New(actualPlatform.Folder).Join("libraries")
		if isDir, _ := actualPlatformLibDir.IsDir(); isDir {
			sortedLibrariesFolders.Add(actualPlatformLibDir)
		}
	}

	platformLibDir := paths.New(platform.Folder).Join("libraries")
	if isDir, _ := platformLibDir.IsDir(); isDir {
		sortedLibrariesFolders.Add(platformLibDir)
	}

	librariesFolders := ctx.OtherLibrariesFolders
	if err := librariesFolders.ToAbs(); err != nil {
		return i18n.WrapError(err)
	}
	sortedLibrariesFolders.AddAllMissing(librariesFolders)

	ctx.LibrariesFolders = sortedLibrariesFolders

	var libs []*libraries.Library
	for _, libraryFolder := range sortedLibrariesFolders {
		libsInFolder, err := libraries.LoadLibrariesFromDir(libraryFolder)
		if err != nil {
			return i18n.WrapError(err)
		}
		libs = append(libs, libsInFolder...)
	}
	if debugLevel > 0 {
		for _, lib := range libs {
			warnings, err := lib.Lint()
			if err != nil {
				return i18n.WrapError(err)
			}
			for _, warning := range warnings {
				logger.Fprintln(os.Stdout, "warn", warning)
			}
		}
	}

	ctx.Libraries = libs

	headerToLibraries := make(map[string][]*libraries.Library)
	for _, library := range libs {
		headers, err := utils.ReadDirFiltered(library.SrcFolder.String(), utils.FilterFilesWithExtensions(".h", ".hpp", ".hh"))
		if err != nil {
			return i18n.WrapError(err)
		}
		for _, header := range headers {
			headerFileName := header.Name()
			headerToLibraries[headerFileName] = append(headerToLibraries[headerFileName], library)
		}
	}

	ctx.HeaderToLibraries = headerToLibraries

	return nil
}
