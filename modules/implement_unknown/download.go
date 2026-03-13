//modules/implement_unknown/download.go
package 
func DownloadModule(url string, dest string) error {
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    out, _ := os.Create(dest)
    defer out.Close()

    _, err = io.Copy(out, resp.Body)
    return err
}