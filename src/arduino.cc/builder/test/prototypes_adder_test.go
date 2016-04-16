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
 * Copyright 2015 Matthijs Kooijman
 */

package test

import (
	"arduino.cc/builder"
	"arduino.cc/builder/types"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPrototypesAdderBridgeExample(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	sketchLocation := filepath.Join("downloaded_libraries", "Bridge", "examples", "Bridge", "Bridge.ino")
	absoluteSketchLocation := strings.Replace(Abs(t, sketchLocation), "\\", "\\\\", -1)

	context := make(map[string]interface{})
	ctx := &types.Context{
		HardwareFolders:         []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"},
		ToolsFolders:            []string{"downloaded_tools"},
		BuiltInLibrariesFolders: []string{"downloaded_libraries"},
		OtherLibrariesFolders:   []string{"libraries"},
		SketchLocation:          sketchLocation,
		FQBN:                    "arduino:avr:leonardo",
		ArduinoAPIVersion:       "10600",
		Verbose:                 true,
	}

	buildPath := SetupBuildPath(t, ctx)
	defer os.RemoveAll(buildPath)

	ctx.DebugLevel = 10

	commands := []types.Command{

		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},

		&builder.ContainerMergeCopySketchFiles{},

		&builder.ContainerFindIncludes{},

		&builder.PrintUsedLibrariesIfVerbose{},
		&builder.WarnAboutArchIncompatibleLibraries{},

		&builder.ContainerAddPrototypes{},
	}

	for _, command := range commands {
		err := command.Run(context, ctx)
		NoError(t, err)
	}

	require.Equal(t, "#include <Arduino.h>\n#line 1 \""+absoluteSketchLocation+"\"\n", ctx.IncludeSection)
	require.Equal(t, "#line 33 \""+absoluteSketchLocation+"\"\nvoid setup();\n#line 46 \""+absoluteSketchLocation+"\"\nvoid loop();\n#line 62 \""+absoluteSketchLocation+"\"\nvoid process(BridgeClient client);\n#line 82 \""+absoluteSketchLocation+"\"\nvoid digitalCommand(BridgeClient client);\n#line 109 \""+absoluteSketchLocation+"\"\nvoid analogCommand(BridgeClient client);\n#line 149 \""+absoluteSketchLocation+"\"\nvoid modeCommand(BridgeClient client);\n#line 33 \""+absoluteSketchLocation+"\"\n", ctx.PrototypesSection)
}

func TestPrototypesAdderSketchWithIfDef(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	context := make(map[string]interface{})
	ctx := &types.Context{
		HardwareFolders:         []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"},
		ToolsFolders:            []string{"downloaded_tools"},
		BuiltInLibrariesFolders: []string{"downloaded_libraries"},
		OtherLibrariesFolders:   []string{"libraries"},
		SketchLocation:          filepath.Join("sketch2", "SketchWithIfDef.ino"),
		FQBN:                    "arduino:avr:leonardo",
		ArduinoAPIVersion:       "10600",
		Verbose:                 true,
	}

	buildPath := SetupBuildPath(t, ctx)
	defer os.RemoveAll(buildPath)

	commands := []types.Command{

		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},

		&builder.ContainerMergeCopySketchFiles{},

		&builder.ContainerFindIncludes{},

		&builder.PrintUsedLibrariesIfVerbose{},
		&builder.WarnAboutArchIncompatibleLibraries{},

		&builder.ContainerAddPrototypes{},
	}

	for _, command := range commands {
		err := command.Run(context, ctx)
		NoError(t, err)
	}

	preprocessed := LoadAndInterpolate(t, filepath.Join("sketch2", "SketchWithIfDef.preprocessed.txt"), context, ctx)
	require.Equal(t, preprocessed, strings.Replace(ctx.Source, "\r\n", "\n", -1))
}

func TestPrototypesAdderBaladuino(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	context := make(map[string]interface{})
	ctx := &types.Context{
		HardwareFolders:         []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"},
		ToolsFolders:            []string{"downloaded_tools"},
		BuiltInLibrariesFolders: []string{"downloaded_libraries"},
		OtherLibrariesFolders:   []string{"libraries"},
		SketchLocation:          filepath.Join("sketch3", "Baladuino.ino"),
		FQBN:                    "arduino:avr:leonardo",
		ArduinoAPIVersion:       "10600",
		Verbose:                 true,
	}

	buildPath := SetupBuildPath(t, ctx)
	defer os.RemoveAll(buildPath)

	commands := []types.Command{

		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},

		&builder.ContainerMergeCopySketchFiles{},

		&builder.ContainerFindIncludes{},

		&builder.PrintUsedLibrariesIfVerbose{},
		&builder.WarnAboutArchIncompatibleLibraries{},

		&builder.ContainerAddPrototypes{},
	}

	for _, command := range commands {
		err := command.Run(context, ctx)
		NoError(t, err)
	}

	preprocessed := LoadAndInterpolate(t, filepath.Join("sketch3", "Baladuino.preprocessed.txt"), context, ctx)
	require.Equal(t, preprocessed, strings.Replace(ctx.Source, "\r\n", "\n", -1))
}

func TestPrototypesAdderCharWithEscapedDoubleQuote(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	context := make(map[string]interface{})
	ctx := &types.Context{
		HardwareFolders:         []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"},
		ToolsFolders:            []string{"downloaded_tools"},
		BuiltInLibrariesFolders: []string{"downloaded_libraries"},
		OtherLibrariesFolders:   []string{"libraries"},
		SketchLocation:          filepath.Join("sketch4", "CharWithEscapedDoubleQuote.ino"),
		FQBN:                    "arduino:avr:leonardo",
		ArduinoAPIVersion:       "10600",
		Verbose:                 true,
	}

	buildPath := SetupBuildPath(t, ctx)
	defer os.RemoveAll(buildPath)

	commands := []types.Command{

		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},

		&builder.ContainerMergeCopySketchFiles{},

		&builder.ContainerFindIncludes{},

		&builder.PrintUsedLibrariesIfVerbose{},
		&builder.WarnAboutArchIncompatibleLibraries{},

		&builder.ContainerAddPrototypes{},
	}

	for _, command := range commands {
		err := command.Run(context, ctx)
		NoError(t, err)
	}

	preprocessed := LoadAndInterpolate(t, filepath.Join("sketch4", "CharWithEscapedDoubleQuote.preprocessed.txt"), context, ctx)
	require.Equal(t, preprocessed, strings.Replace(ctx.Source, "\r\n", "\n", -1))
}

func TestPrototypesAdderIncludeBetweenMultilineComment(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	context := make(map[string]interface{})
	ctx := &types.Context{
		HardwareFolders:         []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"},
		ToolsFolders:            []string{"downloaded_tools"},
		BuiltInLibrariesFolders: []string{"downloaded_libraries"},
		OtherLibrariesFolders:   []string{"libraries"},
		SketchLocation:          filepath.Join("sketch5", "IncludeBetweenMultilineComment.ino"),
		FQBN:                    "arduino:sam:arduino_due_x_dbg",
		ArduinoAPIVersion:       "10600",
		Verbose:                 true,
	}

	buildPath := SetupBuildPath(t, ctx)
	defer os.RemoveAll(buildPath)

	commands := []types.Command{

		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},

		&builder.ContainerMergeCopySketchFiles{},

		&builder.ContainerFindIncludes{},

		&builder.PrintUsedLibrariesIfVerbose{},
		&builder.WarnAboutArchIncompatibleLibraries{},

		&builder.ContainerAddPrototypes{},
	}

	for _, command := range commands {
		err := command.Run(context, ctx)
		NoError(t, err)
	}

	preprocessed := LoadAndInterpolate(t, filepath.Join("sketch5", "IncludeBetweenMultilineComment.preprocessed.txt"), context, ctx)
	require.Equal(t, preprocessed, strings.Replace(ctx.Source, "\r\n", "\n", -1))
}

func TestPrototypesAdderLineContinuations(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	context := make(map[string]interface{})
	ctx := &types.Context{
		HardwareFolders:         []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"},
		ToolsFolders:            []string{"downloaded_tools"},
		BuiltInLibrariesFolders: []string{"downloaded_libraries"},
		OtherLibrariesFolders:   []string{"libraries"},
		SketchLocation:          filepath.Join("sketch6", "/LineContinuations.ino"),
		FQBN:                    "arduino:avr:leonardo",
		ArduinoAPIVersion:       "10600",
		Verbose:                 true,
	}

	buildPath := SetupBuildPath(t, ctx)
	defer os.RemoveAll(buildPath)

	commands := []types.Command{

		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},

		&builder.ContainerMergeCopySketchFiles{},

		&builder.ContainerFindIncludes{},

		&builder.PrintUsedLibrariesIfVerbose{},
		&builder.WarnAboutArchIncompatibleLibraries{},

		&builder.ContainerAddPrototypes{},
	}

	for _, command := range commands {
		err := command.Run(context, ctx)
		NoError(t, err)
	}

	preprocessed := LoadAndInterpolate(t, filepath.Join("sketch6", "LineContinuations.preprocessed.txt"), context, ctx)
	require.Equal(t, preprocessed, strings.Replace(ctx.Source, "\r\n", "\n", -1))
}

func TestPrototypesAdderStringWithComment(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	context := make(map[string]interface{})
	ctx := &types.Context{
		HardwareFolders:         []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"},
		ToolsFolders:            []string{"downloaded_tools"},
		BuiltInLibrariesFolders: []string{"downloaded_libraries"},
		OtherLibrariesFolders:   []string{"libraries"},
		SketchLocation:          filepath.Join("sketch7", "StringWithComment.ino"),
		FQBN:                    "arduino:avr:leonardo",
		ArduinoAPIVersion:       "10600",
		Verbose:                 true,
	}

	buildPath := SetupBuildPath(t, ctx)
	defer os.RemoveAll(buildPath)

	commands := []types.Command{

		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},

		&builder.ContainerMergeCopySketchFiles{},

		&builder.ContainerFindIncludes{},

		&builder.PrintUsedLibrariesIfVerbose{},
		&builder.WarnAboutArchIncompatibleLibraries{},

		&builder.ContainerAddPrototypes{},
	}

	for _, command := range commands {
		err := command.Run(context, ctx)
		NoError(t, err)
	}

	preprocessed := LoadAndInterpolate(t, filepath.Join("sketch7", "StringWithComment.preprocessed.txt"), context, ctx)
	require.Equal(t, preprocessed, strings.Replace(ctx.Source, "\r\n", "\n", -1))
}

func TestPrototypesAdderSketchWithStruct(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	context := make(map[string]interface{})
	ctx := &types.Context{
		HardwareFolders:         []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"},
		ToolsFolders:            []string{"downloaded_tools"},
		BuiltInLibrariesFolders: []string{"downloaded_libraries"},
		OtherLibrariesFolders:   []string{"libraries"},
		SketchLocation:          filepath.Join("sketch8", "SketchWithStruct.ino"),
		FQBN:                    "arduino:avr:leonardo",
		ArduinoAPIVersion:       "10600",
		Verbose:                 true,
	}

	buildPath := SetupBuildPath(t, ctx)
	defer os.RemoveAll(buildPath)

	commands := []types.Command{

		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},

		&builder.ContainerMergeCopySketchFiles{},

		&builder.ContainerFindIncludes{},

		&builder.PrintUsedLibrariesIfVerbose{},
		&builder.WarnAboutArchIncompatibleLibraries{},

		&builder.ContainerAddPrototypes{},
	}

	for _, command := range commands {
		err := command.Run(context, ctx)
		NoError(t, err)
	}

	preprocessed := LoadAndInterpolate(t, filepath.Join("sketch8", "SketchWithStruct.preprocessed.txt"), context, ctx)
	obtained := strings.Replace(ctx.Source, "\r\n", "\n", -1)
	// ctags based preprocessing removes the space after "dostuff", but this is still OK
	// TODO: remove this exception when moving to a more powerful parser
	preprocessed = strings.Replace(preprocessed, "void dostuff (A_NEW_TYPE * bar);", "void dostuff(A_NEW_TYPE * bar);", 1)
	obtained = strings.Replace(obtained, "void dostuff (A_NEW_TYPE * bar);", "void dostuff(A_NEW_TYPE * bar);", 1)
	require.Equal(t, preprocessed, obtained)
}

func TestPrototypesAdderSketchWithConfig(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	sketchLocation := filepath.Join("sketch_with_config", "sketch_with_config.ino")
	absoluteSketchLocation := strings.Replace(Abs(t, sketchLocation), "\\", "\\\\", -1)

	context := make(map[string]interface{})
	ctx := &types.Context{
		HardwareFolders:         []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"},
		ToolsFolders:            []string{"downloaded_tools"},
		BuiltInLibrariesFolders: []string{"downloaded_libraries"},
		OtherLibrariesFolders:   []string{"libraries"},
		SketchLocation:          sketchLocation,
		FQBN:                    "arduino:avr:leonardo",
		ArduinoAPIVersion:       "10600",
		Verbose:                 true,
	}

	buildPath := SetupBuildPath(t, ctx)
	defer os.RemoveAll(buildPath)

	commands := []types.Command{

		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},

		&builder.ContainerMergeCopySketchFiles{},

		&builder.ContainerFindIncludes{},

		&builder.PrintUsedLibrariesIfVerbose{},
		&builder.WarnAboutArchIncompatibleLibraries{},

		&builder.ContainerAddPrototypes{},
	}

	for _, command := range commands {
		err := command.Run(context, ctx)
		NoError(t, err)
	}

	require.Equal(t, "#include <Arduino.h>\n#line 1 \""+absoluteSketchLocation+"\"\n", ctx.IncludeSection)
	require.Equal(t, "#line 13 \""+absoluteSketchLocation+"\"\nvoid setup();\n#line 17 \""+absoluteSketchLocation+"\"\nvoid loop();\n#line 13 \""+absoluteSketchLocation+"\"\n", ctx.PrototypesSection)

	preprocessed := LoadAndInterpolate(t, filepath.Join("sketch_with_config", "sketch_with_config.preprocessed.txt"), context, ctx)
	require.Equal(t, preprocessed, strings.Replace(ctx.Source, "\r\n", "\n", -1))
}

func TestPrototypesAdderSketchNoFunctionsTwoFiles(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	sketchLocation := filepath.Join("sketch_no_functions_two_files", "main.ino")
	absoluteSketchLocation := strings.Replace(Abs(t, sketchLocation), "\\", "\\\\", -1)

	context := make(map[string]interface{})
	ctx := &types.Context{
		HardwareFolders:         []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"},
		ToolsFolders:            []string{"downloaded_tools"},
		BuiltInLibrariesFolders: []string{"downloaded_libraries"},
		OtherLibrariesFolders:   []string{"libraries"},
		SketchLocation:          filepath.Join("sketch_no_functions_two_files", "main.ino"),
		FQBN:                    "arduino:avr:leonardo",
		ArduinoAPIVersion:       "10600",
		Verbose:                 true,
	}

	buildPath := SetupBuildPath(t, ctx)
	defer os.RemoveAll(buildPath)

	commands := []types.Command{

		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},

		&builder.ContainerMergeCopySketchFiles{},

		&builder.ContainerFindIncludes{},

		&builder.PrintUsedLibrariesIfVerbose{},
		&builder.WarnAboutArchIncompatibleLibraries{},

		&builder.ContainerAddPrototypes{},
	}

	for _, command := range commands {
		err := command.Run(context, ctx)
		NoError(t, err)
	}

	require.Equal(t, "#include <Arduino.h>\n#line 1 \""+absoluteSketchLocation+"\"\n", ctx.IncludeSection)
	require.Equal(t, "", ctx.PrototypesSection)
}

func TestPrototypesAdderSketchNoFunctions(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	context := make(map[string]interface{})
	ctx := &types.Context{
		HardwareFolders:         []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"},
		ToolsFolders:            []string{"downloaded_tools"},
		BuiltInLibrariesFolders: []string{"downloaded_libraries"},
		OtherLibrariesFolders:   []string{"libraries"},
		SketchLocation:          filepath.Join("sketch_no_functions", "main.ino"),
		FQBN:                    "arduino:avr:leonardo",
		ArduinoAPIVersion:       "10600",
		Verbose:                 true,
	}

	buildPath := SetupBuildPath(t, ctx)
	defer os.RemoveAll(buildPath)

	sketchLocation := filepath.Join("sketch_no_functions", "main.ino")
	absoluteSketchLocation := strings.Replace(Abs(t, sketchLocation), "\\", "\\\\", -1)

	commands := []types.Command{

		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},

		&builder.ContainerMergeCopySketchFiles{},

		&builder.ContainerFindIncludes{},

		&builder.PrintUsedLibrariesIfVerbose{},
		&builder.WarnAboutArchIncompatibleLibraries{},

		&builder.ContainerAddPrototypes{},
	}

	for _, command := range commands {
		err := command.Run(context, ctx)
		NoError(t, err)
	}

	require.Equal(t, "#include <Arduino.h>\n#line 1 \""+absoluteSketchLocation+"\"\n", ctx.IncludeSection)
	require.Equal(t, "", ctx.PrototypesSection)
}

func TestPrototypesAdderSketchWithDefaultArgs(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	sketchLocation := filepath.Join("sketch_with_default_args", "sketch.ino")
	absoluteSketchLocation := strings.Replace(Abs(t, sketchLocation), "\\", "\\\\", -1)

	context := make(map[string]interface{})
	ctx := &types.Context{
		HardwareFolders:         []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"},
		ToolsFolders:            []string{"downloaded_tools"},
		BuiltInLibrariesFolders: []string{"downloaded_libraries"},
		OtherLibrariesFolders:   []string{"libraries"},
		SketchLocation:          sketchLocation,
		FQBN:                    "arduino:avr:leonardo",
		ArduinoAPIVersion:       "10600",
		Verbose:                 true,
	}

	buildPath := SetupBuildPath(t, ctx)
	defer os.RemoveAll(buildPath)

	commands := []types.Command{

		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},

		&builder.ContainerMergeCopySketchFiles{},

		&builder.ContainerFindIncludes{},

		&builder.PrintUsedLibrariesIfVerbose{},
		&builder.WarnAboutArchIncompatibleLibraries{},

		&builder.ContainerAddPrototypes{},
	}

	for _, command := range commands {
		err := command.Run(context, ctx)
		NoError(t, err)
	}

	require.Equal(t, "#include <Arduino.h>\n#line 1 \""+absoluteSketchLocation+"\"\n", ctx.IncludeSection)
	require.Equal(t, "#line 4 \""+absoluteSketchLocation+"\"\nvoid setup();\n#line 7 \""+absoluteSketchLocation+"\"\nvoid loop();\n#line 1 \""+absoluteSketchLocation+"\"\n", ctx.PrototypesSection)
}

func TestPrototypesAdderSketchWithInlineFunction(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	sketchLocation := filepath.Join("sketch_with_inline_function", "sketch.ino")
	absoluteSketchLocation := strings.Replace(Abs(t, sketchLocation), "\\", "\\\\", -1)

	context := make(map[string]interface{})
	ctx := &types.Context{
		HardwareFolders:         []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"},
		ToolsFolders:            []string{"downloaded_tools"},
		BuiltInLibrariesFolders: []string{"downloaded_libraries"},
		OtherLibrariesFolders:   []string{"libraries"},
		SketchLocation:          sketchLocation,
		FQBN:                    "arduino:avr:leonardo",
		ArduinoAPIVersion:       "10600",
		Verbose:                 true,
	}

	buildPath := SetupBuildPath(t, ctx)
	defer os.RemoveAll(buildPath)

	commands := []types.Command{

		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},

		&builder.ContainerMergeCopySketchFiles{},

		&builder.ContainerFindIncludes{},

		&builder.PrintUsedLibrariesIfVerbose{},
		&builder.WarnAboutArchIncompatibleLibraries{},

		&builder.ContainerAddPrototypes{},
	}

	for _, command := range commands {
		err := command.Run(context, ctx)
		NoError(t, err)
	}

	require.Equal(t, "#include <Arduino.h>\n#line 1 \""+absoluteSketchLocation+"\"\n", ctx.IncludeSection)

	expected := "#line 1 \"" + absoluteSketchLocation + "\"\nvoid setup();\n#line 2 \"" + absoluteSketchLocation + "\"\nvoid loop();\n#line 4 \"" + absoluteSketchLocation + "\"\nshort unsigned int testInt();\n#line 8 \"" + absoluteSketchLocation + "\"\nstatic int8_t testInline();\n#line 12 \"" + absoluteSketchLocation + "\"\n__attribute__((always_inline)) uint8_t testAttribute();\n#line 1 \"" + absoluteSketchLocation + "\"\n"
	obtained := ctx.PrototypesSection
	// ctags based preprocessing removes "inline" but this is still OK
	// TODO: remove this exception when moving to a more powerful parser
	expected = strings.Replace(expected, "static inline int8_t testInline();", "static int8_t testInline();", -1)
	obtained = strings.Replace(obtained, "static inline int8_t testInline();", "static int8_t testInline();", -1)
	// ctags based preprocessing removes "__attribute__ ....." but this is still OK
	// TODO: remove this exception when moving to a more powerful parser
	expected = strings.Replace(expected, "__attribute__((always_inline)) uint8_t testAttribute();", "uint8_t testAttribute();", -1)
	obtained = strings.Replace(obtained, "__attribute__((always_inline)) uint8_t testAttribute();", "uint8_t testAttribute();", -1)
	require.Equal(t, expected, obtained)
}

func TestPrototypesAdderSketchWithFunctionSignatureInsideIFDEF(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	sketchLocation := filepath.Join("sketch_with_function_signature_inside_ifdef", "sketch.ino")
	absoluteSketchLocation := strings.Replace(Abs(t, sketchLocation), "\\", "\\\\", -1)

	context := make(map[string]interface{})
	ctx := &types.Context{
		HardwareFolders:         []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"},
		ToolsFolders:            []string{"downloaded_tools"},
		BuiltInLibrariesFolders: []string{"downloaded_libraries"},
		OtherLibrariesFolders:   []string{"libraries"},
		SketchLocation:          sketchLocation,
		FQBN:                    "arduino:avr:leonardo",
		ArduinoAPIVersion:       "10600",
		Verbose:                 true,
	}

	buildPath := SetupBuildPath(t, ctx)
	defer os.RemoveAll(buildPath)

	commands := []types.Command{

		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},

		&builder.ContainerMergeCopySketchFiles{},

		&builder.ContainerFindIncludes{},

		&builder.PrintUsedLibrariesIfVerbose{},
		&builder.WarnAboutArchIncompatibleLibraries{},

		&builder.ContainerAddPrototypes{},
	}

	for _, command := range commands {
		err := command.Run(context, ctx)
		NoError(t, err)
	}

	require.Equal(t, "#include <Arduino.h>\n#line 1 \""+absoluteSketchLocation+"\"\n", ctx.IncludeSection)
	require.Equal(t, "#line 1 \""+absoluteSketchLocation+"\"\nvoid setup();\n#line 3 \""+absoluteSketchLocation+"\"\nvoid loop();\n#line 15 \""+absoluteSketchLocation+"\"\nint8_t adalight();\n#line 1 \""+absoluteSketchLocation+"\"\n", ctx.PrototypesSection)
}

func TestPrototypesAdderSketchWithUSBCON(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	sketchLocation := filepath.Join("sketch_with_usbcon", "sketch.ino")
	absoluteSketchLocation := strings.Replace(Abs(t, sketchLocation), "\\", "\\\\", -1)

	context := make(map[string]interface{})
	ctx := &types.Context{
		HardwareFolders:         []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"},
		ToolsFolders:            []string{"downloaded_tools"},
		OtherLibrariesFolders:   []string{"libraries"},
		BuiltInLibrariesFolders: []string{"downloaded_libraries"},
		SketchLocation:          sketchLocation,
		FQBN:                    "arduino:avr:leonardo",
		ArduinoAPIVersion:       "10600",
		Verbose:                 true,
	}

	buildPath := SetupBuildPath(t, ctx)
	defer os.RemoveAll(buildPath)

	commands := []types.Command{

		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},

		&builder.ContainerMergeCopySketchFiles{},

		&builder.ContainerFindIncludes{},

		&builder.PrintUsedLibrariesIfVerbose{},
		&builder.WarnAboutArchIncompatibleLibraries{},

		&builder.ContainerAddPrototypes{},
	}

	for _, command := range commands {
		err := command.Run(context, ctx)
		NoError(t, err)
	}

	require.Equal(t, "#include <Arduino.h>\n#line 1 \""+absoluteSketchLocation+"\"\n", ctx.IncludeSection)
	require.Equal(t, "#line 5 \""+absoluteSketchLocation+"\"\nvoid ciao();\n#line 10 \""+absoluteSketchLocation+"\"\nvoid setup();\n#line 15 \""+absoluteSketchLocation+"\"\nvoid loop();\n#line 5 \""+absoluteSketchLocation+"\"\n", ctx.PrototypesSection)
}

func TestPrototypesAdderSketchWithTypename(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	sketchLocation := filepath.Join("sketch_with_typename", "sketch.ino")
	absoluteSketchLocation := strings.Replace(Abs(t, sketchLocation), "\\", "\\\\", -1)

	context := make(map[string]interface{})
	ctx := &types.Context{
		HardwareFolders:   []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"},
		LibrariesFolders:  []string{"libraries", "downloaded_libraries"},
		ToolsFolders:      []string{"downloaded_tools"},
		SketchLocation:    sketchLocation,
		FQBN:              "arduino:avr:leonardo",
		ArduinoAPIVersion: "10600",
		Verbose:           true,
	}

	buildPath := SetupBuildPath(t, ctx)
	defer os.RemoveAll(buildPath)

	commands := []types.Command{

		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},

		&builder.ContainerMergeCopySketchFiles{},

		&builder.ContainerFindIncludes{},

		&builder.PrintUsedLibrariesIfVerbose{},
		&builder.WarnAboutArchIncompatibleLibraries{},

		&builder.ContainerAddPrototypes{},
	}

	for _, command := range commands {
		err := command.Run(context, ctx)
		NoError(t, err)
	}

	require.Equal(t, "#include <Arduino.h>\n#line 1 \""+absoluteSketchLocation+"\"\n", ctx.IncludeSection)
	expected := "#line 6 \"" + absoluteSketchLocation + "\"\nvoid setup();\n#line 10 \"" + absoluteSketchLocation + "\"\nvoid loop();\n#line 12 \"" + absoluteSketchLocation + "\"\ntypename Foo<char>::Bar func();\n#line 6 \"" + absoluteSketchLocation + "\"\n"
	obtained := ctx.PrototypesSection
	// ctags based preprocessing ignores line with typename
	// TODO: remove this exception when moving to a more powerful parser
	expected = strings.Replace(expected, "#line 12 \""+absoluteSketchLocation+"\"\ntypename Foo<char>::Bar func();\n", "", -1)
	obtained = strings.Replace(obtained, "#line 12 \""+absoluteSketchLocation+"\"\ntypename Foo<char>::Bar func();\n", "", -1)
	require.Equal(t, expected, obtained)
}

func TestPrototypesAdderSketchWithIfDef2(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	sketchLocation := filepath.Join("sketch_with_ifdef", "sketch.ino")
	absoluteSketchLocation := strings.Replace(Abs(t, sketchLocation), "\\", "\\\\", -1)

	context := make(map[string]interface{})
	ctx := &types.Context{
		HardwareFolders:         []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"},
		ToolsFolders:            []string{"downloaded_tools"},
		BuiltInLibrariesFolders: []string{"downloaded_libraries"},
		OtherLibrariesFolders:   []string{"libraries"},
		SketchLocation:          sketchLocation,
		FQBN:                    "arduino:avr:yun",
		ArduinoAPIVersion:       "10600",
		Verbose:                 true,
	}

	buildPath := SetupBuildPath(t, ctx)
	defer os.RemoveAll(buildPath)

	commands := []types.Command{

		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},

		&builder.ContainerMergeCopySketchFiles{},

		&builder.ContainerFindIncludes{},

		&builder.PrintUsedLibrariesIfVerbose{},
		&builder.WarnAboutArchIncompatibleLibraries{},

		&builder.ContainerAddPrototypes{},
	}

	for _, command := range commands {
		err := command.Run(context, ctx)
		NoError(t, err)
	}

	require.Equal(t, "#include <Arduino.h>\n#line 1 \""+absoluteSketchLocation+"\"\n", ctx.IncludeSection)
	require.Equal(t, "#line 5 \""+absoluteSketchLocation+"\"\nvoid elseBranch();\n#line 9 \""+absoluteSketchLocation+"\"\nvoid f1();\n#line 10 \""+absoluteSketchLocation+"\"\nvoid f2();\n#line 12 \""+absoluteSketchLocation+"\"\nvoid setup();\n#line 14 \""+absoluteSketchLocation+"\"\nvoid loop();\n#line 5 \""+absoluteSketchLocation+"\"\n", ctx.PrototypesSection)

	expectedSource := LoadAndInterpolate(t, filepath.Join("sketch_with_ifdef", "sketch.preprocessed.txt"), context, ctx)
	require.Equal(t, expectedSource, strings.Replace(ctx.Source, "\r\n", "\n", -1))
}

func TestPrototypesAdderSketchWithIfDef2SAM(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	sketchLocation := filepath.Join("sketch_with_ifdef", "sketch.ino")
	absoluteSketchLocation := strings.Replace(Abs(t, sketchLocation), "\\", "\\\\", -1)

	context := make(map[string]interface{})
	ctx := &types.Context{
		HardwareFolders:         []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"},
		ToolsFolders:            []string{"downloaded_tools"},
		BuiltInLibrariesFolders: []string{"downloaded_libraries"},
		OtherLibrariesFolders:   []string{"libraries"},
		SketchLocation:          sketchLocation,
		FQBN:                    "arduino:sam:arduino_due_x_dbg",
		ArduinoAPIVersion:       "10600",
		Verbose:                 true,
	}

	buildPath := SetupBuildPath(t, ctx)
	defer os.RemoveAll(buildPath)

	commands := []types.Command{

		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},

		&builder.ContainerMergeCopySketchFiles{},

		&builder.ContainerFindIncludes{},

		&builder.PrintUsedLibrariesIfVerbose{},
		&builder.WarnAboutArchIncompatibleLibraries{},

		&builder.ContainerAddPrototypes{},
	}

	for _, command := range commands {
		err := command.Run(context, ctx)
		NoError(t, err)
	}

	require.Equal(t, "#include <Arduino.h>\n#line 1 \""+absoluteSketchLocation+"\"\n", ctx.IncludeSection)
	require.Equal(t, "#line 2 \""+absoluteSketchLocation+"\"\nvoid ifBranch();\n#line 9 \""+absoluteSketchLocation+"\"\nvoid f1();\n#line 10 \""+absoluteSketchLocation+"\"\nvoid f2();\n#line 12 \""+absoluteSketchLocation+"\"\nvoid setup();\n#line 14 \""+absoluteSketchLocation+"\"\nvoid loop();\n#line 2 \""+absoluteSketchLocation+"\"\n", ctx.PrototypesSection)

	expectedSource := LoadAndInterpolate(t, filepath.Join("sketch_with_ifdef", "sketch.preprocessed.SAM.txt"), context, ctx)
	require.Equal(t, expectedSource, strings.Replace(ctx.Source, "\r\n", "\n", -1))
}

func TestPrototypesAdderSketchWithConst(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	sketchLocation := filepath.Join("sketch_with_const", "sketch.ino")
	absoluteSketchLocation := strings.Replace(Abs(t, sketchLocation), "\\", "\\\\", -1)

	context := make(map[string]interface{})
	ctx := &types.Context{
		HardwareFolders:         []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"},
		ToolsFolders:            []string{"downloaded_tools"},
		BuiltInLibrariesFolders: []string{"downloaded_libraries"},
		OtherLibrariesFolders:   []string{"libraries"},
		SketchLocation:          sketchLocation,
		FQBN:                    "arduino:avr:uno",
		ArduinoAPIVersion:       "10600",
		Verbose:                 true,
	}

	buildPath := SetupBuildPath(t, ctx)
	defer os.RemoveAll(buildPath)

	commands := []types.Command{

		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},

		&builder.ContainerMergeCopySketchFiles{},

		&builder.ContainerFindIncludes{},

		&builder.PrintUsedLibrariesIfVerbose{},
		&builder.WarnAboutArchIncompatibleLibraries{},

		&builder.ContainerAddPrototypes{},
	}

	for _, command := range commands {
		err := command.Run(context, ctx)
		NoError(t, err)
	}

	require.Equal(t, "#include <Arduino.h>\n#line 1 \""+absoluteSketchLocation+"\"\n", ctx.IncludeSection)
	require.Equal(t, "#line 1 \""+absoluteSketchLocation+"\"\nvoid setup();\n#line 2 \""+absoluteSketchLocation+"\"\nvoid loop();\n#line 4 \""+absoluteSketchLocation+"\"\nconst __FlashStringHelper* test();\n#line 6 \""+absoluteSketchLocation+"\"\nconst int test3();\n#line 8 \""+absoluteSketchLocation+"\"\nvolatile __FlashStringHelper* test2();\n#line 10 \""+absoluteSketchLocation+"\"\nvolatile int test4();\n#line 1 \""+absoluteSketchLocation+"\"\n", ctx.PrototypesSection)
}

func TestPrototypesAdderSketchWithDosEol(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	context := make(map[string]interface{})
	ctx := &types.Context{
		HardwareFolders:         []string{filepath.Join("..", "hardware"), "hardware", "downloaded_hardware"},
		ToolsFolders:            []string{"downloaded_tools"},
		BuiltInLibrariesFolders: []string{"downloaded_libraries"},
		OtherLibrariesFolders:   []string{"libraries"},
		SketchLocation:          filepath.Join("eol_processing", "sketch.ino"),
		FQBN:                    "arduino:avr:uno",
		ArduinoAPIVersion:       "10600",
		Verbose:                 true,
	}

	buildPath := SetupBuildPath(t, ctx)
	defer os.RemoveAll(buildPath)

	commands := []types.Command{

		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},

		&builder.ContainerMergeCopySketchFiles{},

		&builder.ContainerFindIncludes{},

		&builder.PrintUsedLibrariesIfVerbose{},
		&builder.WarnAboutArchIncompatibleLibraries{},

		&builder.ContainerAddPrototypes{},
	}

	for _, command := range commands {
		err := command.Run(context, ctx)
		NoError(t, err)
	}
	// only requires no error as result
}
