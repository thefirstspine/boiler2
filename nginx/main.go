package nginx

import (
	"fmt"
	"io/ioutil"
)

func GenerateConfig(serverName string, passTo string) string {
	return fmt.Sprintf(`server {
    server_name %s;
    location / {
        proxy_pass http://%s;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
    listen 80;
}`, serverName, passTo)
}

func WriteConfig(filename string, content string) bool {
	byteArray := []byte(content)
	target := fmt.Sprintf("/etc/nginx/sites-enabled/%s", filename)
	fmt.Println("=> write", byteArray, "to", target)
	err := ioutil.WriteFile(target, byteArray, 0644)
	fmt.Println("  => error", err)
	return err == nil
}
