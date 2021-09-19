import sys

def conv_to_gomap(m: dict) -> str:
    temp = ''
    for flag, desc in m.items():
        temp += f'"{flag}": "{desc}",\n'
    return f'map[string]string{{{temp[:-2]}}}'

def conv_to_po(m: dict) -> str:
    temp = ''
    for _, desc in m.items():
        temp += f'msgid "{desc}"\nmsgstr ""\n\n'
    return temp

if __name__ == '__main__':
    with open('flag.txt', 'r') as file:
        m = {}
        lines = file.readlines()
        for line in lines:
            idx = line.index(' ')
            flag = line[:idx].strip()
            desc = line[idx + 1:].strip()
            m[flag] = desc
        
        print(conv_to_gomap(m))
        # print(conv_to_po(m))