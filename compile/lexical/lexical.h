#ifndef _LEXICAL_H_
#define _LEXICAL_H_

#include <string.h>

#ifndef OPERATIONSIGNAL_S
#define OPERATIONSIGNAL_S ",.(){}[]+-*/%=!&|~^><:;?@#`\\\n"
#endif

#ifndef OPERATIONSIGNAL_D

#define OPERATIONSIGNAL_D                                                                                                                              \
    {                                                                                                                                                  \
        "+=", "-=", "*=", "/=", "%=", "==", "!=", "&=", "|=", "~=", "^=", ">=", "<=", ":=", "++", "--", "&&", "||", "=>", "->", "<>", "//", "/*", "*/" \
    }
#endif

const char *movPointer(const char *, const int);

int endNumberConst(const char *, int);

int endStringConst(const char, const char *, int);

int endVariable(const char *, int);

int endOperationSignal(const char *, int);

#endif
