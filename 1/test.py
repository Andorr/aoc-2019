program = [1, 2, 3, 4]

def execute(program):
    instruction = 0
    while True:
        op, a, b, c = program[instruction:instruction + 4]
        if op == 99:
            return
        elif op == 1:
            program[c] = program[a] + program[b]
        elif op == 2:
            program[c] = program[a] * program[b]
        else:
            raise ValueError("invalid program...")
execute(program)
