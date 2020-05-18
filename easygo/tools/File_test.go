/**
 *     ______                 __
 *    /\__  _\               /\ \
 *    \/_/\ \/     ___     __\ \ \         __      ___     ___     __
 *       \ \ \    / ___\ / __ \ \ \  __  / __ \  /  _  \  / ___\ / __ \
 *        \_\ \__/\ \__//\  __/\ \ \_\ \/\ \_\ \_/\ \/\ \/\ \__//\  __/
 *        /\_____\ \____\ \____\\ \____/\ \__/ \_\ \_\ \_\ \____\ \____\
 *        \/_____/\/____/\/____/ \/___/  \/__/\/_/\/_/\/_/\/____/\/____/
 *
 *
 *                                                                    @寒冰
 *                                                            www.icezzz.cn
 *                                                     hanbin020706@163.com
 */
package tools

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestCreateFile(t *testing.T) {

	file, _ := CreateFile("./aaa/bbb.json")
	file.Write([]byte("123123123123"))
	file.Close()
	fileTime := time.Now().Format("2006-01-02")
	newPath := fmt.Sprintf("./aaa/errors_%s.log", fileTime)
	_ = os.Rename("./aaa/bbb.json", newPath)
	file, _ = os.OpenFile("./aaa/bbb.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	oldFile, _ := os.OpenFile(newPath, os.O_WRONLY, 0644)
	_ = Compress([]*os.File{oldFile}, fmt.Sprintf("%s_debug.tar.gz", fileTime))
	oldFile.Close()
	_ = os.Remove(newPath)
	file.Write([]byte("sdfsdfsdfsdfsdf"))
	file.Close()
}
