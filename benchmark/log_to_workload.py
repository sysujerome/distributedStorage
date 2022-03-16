#! /usr/bin/python
import sys

OPS = ["INSERT", "READ", "UPDATE", "DELETE"]

if (len(sys.argv) < 2):
    print("usage : {} input_file".format(sys.argv[0]))  

input_filename = sys.argv[1]
output_filename = input_filename+".formated"
print("write in " + output_filename)


input_file = open(input_filename)
input_lines = input_file.readlines()

output_file = open(output_filename, "w")
output_lines = []

for l in input_lines:
    ls = l.split()
    if (len(ls) < 1):
        continue
    if ls[0] in OPS:
        line = "{} {}\n".format(ls[0], ls[2][4:])
        output_lines.append(line)
        #print(line)
    elif ls[0] == "SCAN":
        line = "{} {} {}\n".format(ls[0], ls[2][4:], ls[3])
        output_lines.append(line)


print(len(output_lines) , " OPs")
output_file.writelines(output_lines)

input_file.close()
output_file.close()