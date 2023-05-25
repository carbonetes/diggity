package files

import (
	"fmt"
	"os"
	"testing"
)

func TestCreateDirectoryWithPermissions(t *testing.T) {
	dirpath := "C:\\Users\\sairen\\AppData\\Local\\Temp\\1139655944"
	// Clean up directory if it already exists
	defer os.RemoveAll(dirpath)

	// Create directory with 0777 permissions
	err := CreateDirectoryWithPermissions(dirpath)
	if err != nil {
		t.Fatalf("Failed to create directory: %s", err)
	}

	// Check if directory was created with correct permissions
	info, err := os.Stat(dirpath)
	if err != nil {
		t.Fatalf("Failed to get directory info: %s", err)
	}

	perms := info.Mode().Perm()
	if perms != 0777 {
		t.Fatalf("Directory permissions are incorrect. Expected 0777, got %o", perms)
	}
}

func CreateDirectoryWithPermissions(dirpath string) error {

	// Check if the directory exists
	if _, err := os.Stat(dirpath); os.IsNotExist(err) {
		// Directory does not exist, create it with 0777 permissions
		err := os.MkdirAll(dirpath, 0777)
		if err != nil {
			fmt.Println(err)
			return err
		}
		defer os.RemoveAll(dirpath)
	} else {
		// Directory already exists, modify its permissions to 0777
		err := os.Chmod(dirpath, 0777)
		fmt.Print(dirpath)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}
