import os
import subprocess
import sys


class GenProto:
    goOutPath = r"gen\v1"

    def __init__(self, workDir):
        self.workDir = workDir

    def check(self):
        g = os.walk(self.workDir)
        for path, dir_list, file_list in g:
            for file_name in file_list:
                abPath = os.path.join(path, file_name)
                if abPath.endswith(".proto"):
                    filePrefix=file_name.split('.')[0]
                    self.done(path, file_name, filePrefix)


    def makedirs(self, path):
        try:
            os.stat(path)
        except:
            os.makedirs(path)

    # 运行命令
    def done(self, path, filename, filePrefix):
        # 当前 proto 文件所在目录为起点 建立 gen/v1
        pre = os.getcwd()

        os.chdir(path)
        goOutPath = os.path.join(self.goOutPath, filePrefix)
        self.makedirs(goOutPath)
        # 当前目录下
        goOutPath = '.'
        cmd = []
        cmd.append("protoc   --go_out=paths=source_relative:%s %s"%(goOutPath, filename))
        cmd.append("protoc   --go-grpc_out=paths=source_relative:%s %s"%(goOutPath, filename))


        #cmd.append( "protoc   --grpc-gateway_out=paths=source_relative,grpc_api_configuration=%s.yaml:%s %s"%(filePrefix, goOutPath,
        for c in cmd:
            try:
                subprocess.run(c, shell=True, check=True)
            except:
                print("some thing error")
                subprocess.run("exit 1", shell=True)

        print("done ", filename)
        os.chdir(pre)

if __name__ == '__main__':
    a = GenProto(".")
    a.check()
