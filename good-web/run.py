import os
s = "api config forms global init middlewares proto router utils validator"
s = s.split(' ')
for i in s:
    try:
        os.mkdir(i)
    except Exception as e:
        print(e.__str__())