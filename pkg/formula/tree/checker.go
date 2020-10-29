package tree

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

// TODO rit home dir must be a parameter
type CheckerManager struct {
	dir  stream.DirLister
	file stream.FileReader
}

func NewChecker(dir stream.DirLister, file stream.FileReader) CheckerManager {
	return CheckerManager{dir: dir, file: file}
}

// CheckCommands is used to warn the user about conflicting
// formula commands on different repos. This function don't
// return err because print an error because of a unsuccessful
// warning attempt can be confusing to the user.
func (cm CheckerManager) CheckCommands() {
	trees := cm.readCommands()
	commands := filterCommands(trees)
	conflictingCommands := conflictingCommands(commands)
	printConflictingCommandsWarning(conflictingCommands)

}

func (cm CheckerManager) readCommands() []formula.Tree {
	repoDir := filepath.Join(api.RitchieHomeDir(), "repos")
	repos, _ := cm.dir.List(repoDir, false)

	tree := formula.Tree{}
	var treeArr []formula.Tree
	for _, r := range repos {
		path := fmt.Sprintf(treeRepoCmdPattern, api.RitchieHomeDir(), r)
		bytes, _ := cm.file.Read(path)
		_ = json.Unmarshal(bytes, &tree)
		treeArr = append(treeArr, tree)
		tree = formula.Tree{}
	}

	return treeArr
}

func filterCommands(tree []formula.Tree) []string {
	allCommands := []string{""}
	for _, t := range tree {
		for _, c := range t.Commands {
			allCommands = append(allCommands, c.Id)
		}
	}
	return allCommands
}

func conflictingCommands(commands []string) []string {
	duplicateFrequency := make(map[string]int)
	duplicatedCommands := []string{""}
	for _, item := range commands {
		_, exist := duplicateFrequency[item]
		if exist {
			duplicateFrequency[item] += 1 // increase counter by 1 if already in the map
			duplicatedCommands = append(duplicatedCommands, item)
		} else {
			duplicateFrequency[item] = 1 // else start counting from 1
		}
	}
	return duplicatedCommands

}


func printConflictingCommandsWarning(conflictingCommands []string) {
	fmt.Print(prompt.Yellow("The following formula commands are conflicting:"))
	fc := formatCommands(conflictingCommands)
	for _, c := range fc {
		fmt.Println(c)
	}
}

func formatCommands(commands []string)  []string {
	for i := range commands {
        commands[i] = strings.Replace(commands[i], "root", "rit", 1)
		commands[i] = strings.ReplaceAll(commands[i], "_", " ")
	}
	return commands
}
