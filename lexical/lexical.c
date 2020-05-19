#include <stdio.h>
#include "lexical.h"

int endNumberConst(const char *str, int length)
{
    int i;

    if (!length)
    {
        length = strlen(str);
    }

    /* 10进制 */
    if ('0' != str[0])
    {
        for (i = 0; i < length && (('0' <= str[i] && str[i] <= '9') || '.' == str[i]); i++)
        {
        }
        return i;
    }

    /* hex */
    if ('x' == str[1] || 'X' == str[1])
    {
        for (i = 2; i < length && (('0' <= str[i] && str[i] <= '9') || ('A' <= str[i] && str[i] <= 'F')); i++)
        {
        }
        return i;
    }

    /* oct */
    for (i = 0; i < length && ('0' <= str[i] && str[i] < '8'); i++)
    {
    }
    return i;
}

int endStringConst(const char endCh, const char *str, int length)
{
    int i;

    if (!length)
    {
        length = strlen(str);
    }

    for (i = 1; i < length && endCh != str[i]; i++)
    {
        if ('\\' == str[i] && (i + 1) < length)
        {
            i++;
        }
    }

    return i;
}

int endVariable(const char *str, int length)
{
    int i;

    if (!length)
    {
        length = strlen(str);
    }

    for (
        i = 0;
        i < length &&
        (
            ('A' <= str[i] && str[i] <= 'Z') ||
            ('a' <= str[i] && str[i] <= 'z') ||
            ('0' <= str[i] && str[i] <= '9') ||
            '$' == str[i] ||
            '_' == str[i] ||
            0x80 <= (unsigned char)str[i]
        );
        i++)
    {
    }

    return i;
}

int endOperationSignal(const char *str, int length)
{
    char ch;
    short wchar, *wp;
    const char lst_s[] = OPERATIONSIGNAL_S;
    const char lst_d[][3] = OPERATIONSIGNAL_D;
    register int i, len;

    len = (sizeof lst_d) / 3;

    if (!length)
    {
        length = strlen(str);
    }

    if (1 < length)
    {
        wchar = *((const short *)str);
        for (i = 0; i < len; i++)
        {
            wp = (short *)lst_d[i];
            if (wchar == *wp)
            {
                return 2;
            }
        }
    }

    ch = *str;
    for (i = 0; i < sizeof lst_s; i++)
    {
        if (ch == lst_s[i])
        {
            return 1;
        }
    }

    return 0;
}

const char *movPointer(const char *str, const int offset)
{
    return str + offset;
}
