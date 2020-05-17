# Timing Attack

## Intro

A [timing attack](https://en.wikipedia.org/wiki/Timing_attack) is a type of side channel attack (meaning it exploits a vulnerability in the implementation of the sercurity system rather than an inherent weakness of that system) where the time taken to execute a cryptographic algorithm is used to gain information about and eventually compromise the system. A system is vulnerable to timing attacks if the time taken to execute the algorithm varies with the degree to which the attacker's attempt is correct.

In this example, we will implement a simple server that accepts a password and responds with information if the password is correct and with an error if incorrect. By measuring the time taken to verify the given password we will attempt to crack the password.

The benefit of using this attack over a brute-force method has to do with timing. If we say that a password can use `n` different characters and is of length `m`, then the time to perform a brute-force attack (assming you know the value of `m`) is `n^m`, while the timing attack method takes only `n * m * k` time, where `k` is the number of samples to take per attempted password. So with 26 characters available and a password length of 10, the brute-force method has a time cost of 26^10, or 141,167,095,653,376, while the timing attack approach, with a sample size of 10,000, has a cost of only 26 * 10 * 10,000, or 2,600,000—faster by a factor of about 54 million. If, for example, the timing attack took a week, the brute-force method would take about 1 million years.

Note that increasing the sample size increase the time needed to perform the attack, but makes timing estimates more reliable.

## Vulnerability

To check if two passwords match (the given and the real passwords) we can iterate over the characters in each of the passwords and compare them to each other. As soon as a single character does not match, we know that the passwords do not match and can respond to the client telling them that the password they have entered is incorrect. This means that the greater the number of correct characters at the beginning of the attempted password, the longer the check takes—and that gives us enough information to crack the password.

*Note*: I tried to break Go's built-in string comparison operator (`==` or `!=`) and was unable to find reliable timing differences between different passwords, so I've instead opted to go with a more explicit loop implementation of string checking (with a delay to exaggerate the timing). While a little contrived, I hope it illustrates the point.

*Note*: Also in the interest of time, I've made the password short and the character set small, but these could be increased and handled by increasing the amount of sampling for each attempted password.

To determine which password is most correct out of a set of options (with timing measurements), the attacker can use statistical methods to compare averages. In this example, we'll just try each option an equal number of times and then select the password whose total time is greatest. A few improvements could be made on this by making use of other statistical/search/feedback methods like [multi-armed bandit](https://en.wikipedia.org/wiki/Multi-armed_bandit) solutions or [genetic algorithms](https://en.wikipedia.org/wiki/Genetic_algorithm).

## Solution

To protect against this vulnerability, the server should find a way to compare passwords in a constant-time fashion, so that all password attempts take the same amount of time to check on average. In this example, the server checks all characters in the passwords even after finding mismatches and only then respond to the client. Another approach would be to hash the passwords so that all attempts have the same length and compare the hashes.

## Running the program(s)

To run the server, run:

```shell
go run src/server/server.go
```

To run the attacker, run:

```shell
go run src/attacker/attacker.go [-v] [-s]
```

The verbose (`v`) flag controls logging output, setting it to true (including the flag) will cause character timing data to be printed for each character position. The secure (`s`) flag sets which endpoint to try—either the secure or vulnerable endpoint—setting it to true causes the attacker to attempt the secure endpoint.