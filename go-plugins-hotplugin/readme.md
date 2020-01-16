## Go Plugins using hotplugin

Testing [hotplugin](https://github.com/letiantech/hotplugin).

### Running

1. Build and deploy the plugin.
   ```bash
   $ cd plugins/testplugin && ./build-deploy.sh && cd -
   ```
2. Run the app.
   ```bash
   $ go run main.go
   ```

Result:

```bash
go-plugins-hotplugin > go run main.go
2020/01/15 14:45:50 ./plugins/testplugin.so
2020/01/15 14:45:50 load plugin: TestPlugin, version: 0x10000
2020/01/15 14:45:50 TestPlugin loaded
2020/01/15 14:45:50 TestPlugin > loaded version: 0x10000
hello my world
go-plugins-hotplugin >
```
