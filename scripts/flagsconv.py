import sys

def conv_to_gomap(lines: list[str]) -> str:
    temp = ''
    for line in lines:
        parts = line.split('//')
        flag = parts[0].strip().removeprefix('MF_').removesuffix(',')
        desc = parts[1].strip()
        temp += f'"{flag}": "{desc}",\n'
    return f'map[string]string{{{temp[:-2]}}}'

def conv_to_po(lines: list[str]) -> str:
    res = ''
    for line in lines:
        parts = line.split('//')
        desc = parts[1].strip()
        res += f'msgid "{desc}"\nmsgstr ""\n\n'
    return res

if __name__ == '__main__':
    with open('flag.txt', 'r') as file:
        lines = file.readlines()
        print(conv_to_gomap(lines))
        print(conv_to_po(lines))