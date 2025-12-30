package lessons

import "github.com/Bparsons0904/rootcamp/internal/types"

func GetAll() []types.Lesson {
	return []types.Lesson{
		{
			ID:    "cd",
			Title: "Navigate Directories",
			Order: 1,
			What: `The 'cd' command changes your current directory. Think of it as moving between folders - you're always "standing" in one directory at a time, and cd lets you move to a different one.`,
			Example: `cd Documents
cd /home/user/projects
cd ..              # Go up one directory
cd ~               # Go to home directory
cd -               # Go to previous directory`,
			History: `The name 'cd' stands for "change directory". It's been a fundamental command since the early Unix systems in the 1970s. The concept of hierarchical directories was a major innovation that made file organization intuitive and scalable.`,
			CommonUses: []string{
				"cd ~/Documents - Navigate to your Documents folder",
				"cd .. - Move up one level in the directory tree",
				"cd /var/log - Go to a specific system directory",
				"cd ~ - Return to your home directory from anywhere",
				"cd - - Toggle between current and previous directory",
			},
			Lab: types.LabConfig{
				Dirs: []string{"projects", "projects/web", "projects/cli", "documents"},
				Files: map[string]string{
					"projects/web/secret.txt": "CD_MASTER",
					"README.txt":              "Welcome to the lab! Navigate to projects/web/ and find the secret code.",
				},
				Task: "Navigate to the 'projects/web' directory and read the secret.txt file to find your completion code. Hint: You'll need to use 'cd' to get there and 'cat' to read the file.",
			},
			SecretCode: "CD_MASTER",
			Difficulty: "beginner",
			Tags:       []string{"navigation", "basics", "filesystem"},
		},
		{
			ID:    "ls",
			Title: "List Directory Contents",
			Order: 2,
			What: `The 'ls' command lists files and directories in your current location. It's like opening a folder to see what's inside, but in the terminal.`,
			Example: `ls                 # List current directory
ls -l              # Long format with details
ls -a              # Show hidden files (starting with .)
ls -lh             # Human-readable file sizes
ls /etc            # List contents of /etc`,
			History: `'ls' stands for "list" and has been part of Unix since the beginning. The command was designed to be short and quick to type - a philosophy that influenced many Unix command names. The various flags (-l, -a, etc.) were added over time as users needed more detailed information.`,
			CommonUses: []string{
				"ls -la - Show all files with detailed information",
				"ls -lh - Display file sizes in human-readable format (KB, MB, GB)",
				"ls -t - Sort files by modification time (newest first)",
				"ls -R - Recursively list all subdirectories",
				"ls *.txt - List only text files in current directory",
			},
			Lab: types.LabConfig{
				Dirs: []string{"workspace", "workspace/hidden"},
				Files: map[string]string{
					"workspace/file1.txt":       "Regular file",
					"workspace/file2.txt":       "Another file",
					"workspace/.hidden_file":    "I'm hidden! Use ls -a to see me.",
					"workspace/.secret_code":    "LS_EXPLORER",
					"workspace/hidden/dummy.txt": "Not the code",
				},
				Task: "Use 'ls' with the appropriate flag to find the hidden file that contains your secret code. Remember: hidden files start with a dot (.).",
			},
			SecretCode: "LS_EXPLORER",
			Difficulty: "beginner",
			Tags:       []string{"navigation", "basics", "filesystem"},
		},
		{
			ID:    "cat",
			Title: "Display File Contents",
			Order: 3,
			What: `The 'cat' command displays the contents of a file in your terminal. It's one of the most common ways to quickly read text files without opening an editor.`,
			Example: `cat file.txt                    # Display file contents
cat file1.txt file2.txt         # Display multiple files
cat file.txt | grep "search"    # Combine with other commands
cat > newfile.txt               # Create new file (Ctrl+D to save)`,
			History: `'cat' is short for "concatenate" - its original purpose was to concatenate multiple files together. The name might seem odd for just viewing files, but concatenation is still one of its key features. Created by Ken Thompson as part of the original Unix in 1971.`,
			CommonUses: []string{
				"cat README.md - Quickly view a file's contents",
				"cat error.log - Check log files for errors",
				"cat file1.txt file2.txt > combined.txt - Merge files",
				"cat /etc/hosts - View system configuration files",
				"cat << EOF > file.txt - Create file with multi-line content",
			},
			Lab: types.LabConfig{
				Dirs: []string{"data"},
				Files: map[string]string{
					"data/note1.txt": "This is not the code you're looking for.",
					"data/note2.txt": "Keep searching!",
					"data/note3.txt": "The code is: CAT_READER\n\nCongratulations on finding it!",
					"hint.txt":       "Try reading each file in the data directory using cat.",
				},
				Task: "Use 'cat' to read files in the data directory until you find the one containing the secret code. You may need to use 'cd' and 'ls' first!",
			},
			SecretCode: "CAT_READER",
			Difficulty: "beginner",
			Tags:       []string{"file-reading", "basics", "text"},
		},
		{
			ID:    "cp",
			Title: "Copy Files and Directories",
			Order: 4,
			What: `The 'cp' command copies files and directories. It creates an exact duplicate in a new location while keeping the original intact.`,
			Example: `cp source.txt dest.txt          # Copy file
cp source.txt /path/to/dest/    # Copy to directory
cp -r folder/ backup/           # Copy directory recursively
cp *.txt documents/             # Copy all .txt files`,
			History: `'cp' stands for "copy" and was part of early Unix systems. Before graphical interfaces, copying files via command line was the only option. The '-r' (recursive) flag was added later to handle directory copying, as early versions only worked with individual files.`,
			CommonUses: []string{
				"cp config.yml config.yml.backup - Create backup of important file",
				"cp -r project/ project-backup/ - Backup entire directory",
				"cp /etc/config.conf . - Copy system file to current directory",
				"cp -i file.txt dest.txt - Interactive mode (prompt before overwriting)",
				"cp -u source.txt dest.txt - Copy only if source is newer",
			},
			Lab: types.LabConfig{
				Dirs: []string{"original", "backup"},
				Files: map[string]string{
					"original/important.txt": "CP_DUPLICATOR",
					"original/data.txt":      "Some data",
					"instructions.txt":       "Copy the important.txt file from the original/ directory to the backup/ directory, then read it from the backup location.",
				},
				Task: "Use 'cp' to copy important.txt from the original/ directory to the backup/ directory. Then use 'cat' to read the copied file and get your secret code.",
			},
			SecretCode: "CP_DUPLICATOR",
			Difficulty: "beginner",
			Tags:       []string{"file-management", "basics", "copying"},
		},
		{
			ID:    "mv",
			Title: "Move and Rename Files",
			Order: 5,
			What: `The 'mv' command moves files to a new location or renames them. Unlike 'cp', it doesn't create a duplicate - it relocates the original file.`,
			Example: `mv oldname.txt newname.txt      # Rename file
mv file.txt /path/to/dest/      # Move file
mv *.jpg photos/                # Move all .jpg files
mv -i file.txt dest.txt         # Interactive (confirm before overwrite)`,
			History: `'mv' stands for "move" and has been fundamental to Unix since the 1970s. Interestingly, renaming is just a special case of moving - you're moving the file to the same directory with a different name. The command handles both operations because at the filesystem level, they're the same thing.`,
			CommonUses: []string{
				"mv old_name.txt new_name.txt - Rename a file",
				"mv file.txt ~/Documents/ - Move file to Documents folder",
				"mv *.log logs/ - Organize files by moving to subdirectory",
				"mv draft.txt final.txt - Rename file when project is complete",
				"mv -n source.txt dest.txt - Don't overwrite existing files",
			},
			Lab: types.LabConfig{
				Dirs: []string{"messy", "organized"},
				Files: map[string]string{
					"messy/wrong_name.txt": "MV_ORGANIZER",
					"messy/file1.txt":      "Random file",
					"messy/file2.txt":      "Another file",
					"task.txt":             "There's a file called 'wrong_name.txt' in the messy/ directory. Rename it to 'secret.txt' using mv, then read it.",
				},
				Task: "Use 'mv' to rename 'wrong_name.txt' to 'secret.txt' in the messy/ directory. Then use 'cat' to read the renamed file and find your code.",
			},
			SecretCode: "MV_ORGANIZER",
			Difficulty: "beginner",
			Tags:       []string{"file-management", "basics", "organizing"},
		},
	}
}

func GetByID(id string) *types.Lesson {
	lessons := GetAll()
	for i := range lessons {
		if lessons[i].ID == id {
			return &lessons[i]
		}
	}
	return nil
}
