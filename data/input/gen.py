#!/usr/bin/env python3
import sys
import random
from collections import Counter


def generate_files(num_lines, file_number):
    # List of 32 sample words
    word_list = [
        "algorithm",
        "computer",
        "database",
        "network",
        "software",
        "hardware",
        "internet",
        "program",
        "system",
        "memory",
        "processor",
        "security",
        "interface",
        "protocol",
        "server",
        "client",
        "browser",
        "cache",
        "kernel",
        "compiler",
        "debug",
        "encrypt",
        "firewall",
        "gateway",
        "hosting",
        "index",
        "java",
        "kernel",
        "linux",
        "module",
        "node",
        "python",
    ]

    try:
        # Convert input argument to integer
        num_lines = int(num_lines)
        file_number = int(file_number)

        if num_lines <= 0:
            raise ValueError("Number of lines must be positive")

        input_filename = f"input_{file_number}.txt"
        output_filename = f"output_{file_number}.txt"

        # Keep track of all words for counting
        all_words = []

        # Generate input file
        with open(input_filename, "w") as f:
            # Write the number of lines as first line
            f.write(f"{num_lines}\n")

            # Generate specified number of lines
            for _ in range(num_lines):
                # Randomly choose how many words for this line (1 to 15)
                num_words = random.randint(1, 15)

                # Generate the line by sampling words
                line_words = random.choices(word_list, k=num_words)
                all_words.extend(line_words)

                # Write the line
                f.write(" ".join(line_words) + "\n")

        # Generate output file with word counts
        word_counts = Counter(all_words)
        with open(output_filename, "w") as f:
            for word, count in sorted(word_counts.items()):
                f.write(f"{word}: {count}\n")

        print(f"Successfully generated:")
        print(f"1. {input_filename} with {num_lines} lines")
        print(f"2. {output_filename} with word counts")

    except ValueError as e:
        print(f"Error: {e}")
        sys.exit(1)


if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("Usage: python script.py <number_of_lines> <file_number>")
        print("Example: python script.py 78 1")
        sys.exit(1)

    generate_files(sys.argv[1], sys.argv[2])
