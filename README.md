# 华为音乐下载工具

本项目是一个华为音乐下载工具，提供了多个参数用于控制下载行为。该工具支持单曲、专辑、歌单的下载，并允许用户自定义文件输出路径、下载质量、并发数量等。
## 版权声明
本工具仅供学习和研究使用，请勿用于商业用途或违反法律规定。

## 使用示例
###  基本用法
**先编译**
```sh
go build -o hwyy #只需要编译一次
```
**直接双击**可以手动输入url，会自动识别是单曲、歌单或者其他

**命令行调用**
```sh
./hwyy [flags] url
```
**查看帮助**
```sh
./hwyy -h 
```
单曲 sq
```sh
./hwyy -o ./output -q 3   https://url.cloud.huawei.com/rNNkOMZHB6
```
- `-o ./output` 指定输出目录
- `-q 3` 选择 SQ（无损品质）

下载歌单 前三首 空间音频
```sh
./hwyy  -q 13 -r 1-3 https://url.cloud.huawei.com/qdnlrLSU24
```


## 配置文件
配置文件名：[config.yaml](config.yaml)`config.yaml`
## 参数说明
- **authorization**: 需要抓包，每次重启 APP 都会更换，建议使用模拟器
- **输出文件命名格式**
    ```yaml
    file_name: '{title}-{filesize}-{duration}-{rate} '
    ```
- **单曲输出路径**
  ```yaml
  single_format: '{output}/{artist}/'
  ```
- **专辑输出路径**
  ```yaml
  album_format: '{output}/{album}/'
  ```
- **歌单输出格式**
  ```yaml
  playlist_format: '{output}/{playlist_name}/'
  ```
- **歌手单曲输出路径**
  ```yaml
  artist_single_format: '{output}/{artist}/'
  ```
- **歌手专辑输出路径**
  ```yaml
  artist_album_format: '{output}/{artist}/{album}/'
  ```
- **可选变量**

| 变量名           | 说明    |
|---------------|-------|
| artist        | 歌手    |
| title         | 标题    |
| album         | 专辑名   |
| sub_title     | 副标题   |
| output        | 输出路径  |
| playlist_name | 歌单名   |
| fileFormat    | 音频格式  |
| filesize      | 文件大小  |
| duration      | 音频时长  |
| rate          | 采样率   |
| file_name     | 文件名格式 |

## 可选设置(同时支持命令行和配置文件)
- **是否下载歌词 (-l 参数)**
  ```yaml
  lyric: true
  ```
- **是否下载封面 (-c 参数)**
  ```yaml
  cover: true
  ```
- **封面大小选择 (-cv 参数)**
  ```yaml
  cover_size: big  
  # 可选: big (1000*1000), mid (600*600), small (320*320)
  ```
- **输出路径 (-o 参数)**
  ```yaml
  output: ./o
  ```
- **单次批量解析的最大数量 (-m 参数)**
  ```yaml
  max_count: 500  # 默认为100
  ```
- **歌手下载类型 (-a 参数)**
  ```yaml
  artist_type: s  
  #s: 单曲a: 专辑, 默认 s
  ```
- **音质选择 (-q 参数)**
  ```yaml
   quality: 1  
  1 标准品质
  2 HQ
  3 SQ
  4 HIFI
  5 Hi-Res
  13 Audio Vivid
  15 多轨道
  all 全部
  best 最大体积
  具体情况具体判断，不全
  ```
  如果指定 `quality`不存在，下载时会提示手动选择。

- **歌单/专辑/歌手/搜索单曲的下载范围 (-r 参数)**
  ```yaml
  range: all  # 默认 1-3，可指定具体范围，如 1-10,13,20-30
  ```
- **下载歌手专辑时的专辑范围 (-z 参数)**
  ```yaml
  album_range: 1-10,13,20-30  # 可选，默认 1-3
  ```
- **多线程下载数量 (-d 参数)**
  ```yaml
  num_threads: 10  # 默认为5
  ```
  
如果同时使用**配置文件**和**命令行参数**，**命令行参数优先**。