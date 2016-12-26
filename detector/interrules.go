package detector

var interalRules string = `- language: php
  top:
  - score: 100
    value: ^<\?php
  - score: 100
    value: '\#!/usr/bin/env\s+php'
  neartop:
  - score: 2
    value: use( )+\w+(\\\w+)+( )*;
  - score: 3
    value: ^<\?
  - score: -5
    value: ^<\%
  rules:
  - score: 2
    value: <\?php
  - score: 2
    value: ^\?>
  - score: 2
    value: \$\w+
  - score: 2
    value: \$\w+\->\w+
  - score: 2
    value: (require|include)(_once)?( )*\(?( )*('|").+\.php('|")( )*\)?( )*;
  - score: 1
    value: echo( )+('|").+('|")( )*;
  - score: 1
    value: 'NULL'
  - score: 1
    value: new( )+((\\\w+)+|\w+)(\(.*\))?
  - score: 1
    value: function(( )+[\$\w]+\(.*\)|( )*\(.*\))
##asp use funcion ... end function	
  - score: -1
    value: (?i)\bend\s*function\b
  - score: 1
    value: (else)?if( )+\(.+\)
  - score: 1
    value: \w+::\w+
  - score: 1
    value: ===
  - score: 1
    value: '!=='
  - score: -1
    value: (^|\s)(var|char|long|int|float|double)( )+\w+( )*=?
  - score: -5
    value: (?i)\b(dim|next)\b			
  - score: -2
    value: ^\%[<>]
  - score: -10
    value: '(?i)^<%@\s+import\s+Namespace="?System.*"?\s*%>'
- language: jsp
  top:
  - score: 100
    value: '^<%@\s*page\s+language="java.*%>'
  - score: 100
    value: '<%@\s*page\s*import="java.*".*%>'
  neartop:
  - score: 30
    value: '<%@\s*page\s+language="java'
  - score: 30
    value: '<%@\s*page\s*.*import="java\.'
  - score: 10
    value: \bjava\.io
  rules:
  - score: 2
    value: \}\s*finally\s*\{		
  - score: 2
    value: \bString\b
  - score: 5
    value: '(Runtime\.getRuntime|request\.getParameter)'
  - score: -5
    value: (?i)\s+(dim|next)\s+
  - score: -5
    value: ^<\?
  - score: -5
    value: \?$
  - score: 2
    value: <%.*(request\.getParameter|java\.io).*%>
  - score: 1
    value: new( )+((\\\w+)+|\w+)(\(.*\))?
- language: asp
  top:
  - score: 100
    value: (?i)^<%@\s*language\s*=\s*"?vbscript"?\s*.*%>
  - score: 100
    value: (?i)<!--.*include\s+(file|VIRTUAL)=".*\.asp"\s*-->
  neartop:
  - score: 10
    value: (?i)<%@\s*LANGUAGE\s*=\s*VBScript\.Encode\s*.*%>
  - score: 10
    value: (?i)<SCRIPT LANGUAGE="?VBScript"[a-zA-Z0-9\s]*>
  - score: 10
    value: (?i)<!--.*include\s+(file|VIRTUAL)=".*\.asp"\s*-->
  - score: 4
    value: (?i)<script.*runat="?server"?	
  rules:
  - score: 10
    value: (?i)\bServer\.(Execute|ScriptTimeOut|createobject)\s*
  - score: 10
    value: (?i)(adodb\.connection|ADODB\.RecordSet|ReDim|ADODB\.Stream)
  - score: 10
    value: (?i)(Application|server\.urlEncode)\("?[A-Za-z_0-9]+"?\)
  - score: 5
    value: (?i)(objLinks\.item|CreateTextFile|trace|session\.Contents|CreateNavLink|WshShell\.Run)
  - score: 5
    value: (?i)^[<>]\%
  - score: 1
    value: new( )+((\\\w+)+|\w+)(\(.*\))?
  - score: 1
    value: (?i)\b(request|eval)
  - score: 3
    value: (?i)\bif\b.*\bthen\b	
  - score: 2
    value: (?i)\brequest\.(querystring|form|ServerVariables|BinaryWrite|ContentType|Charset|AddHeader)
  - score: 2
    value: (?i)Response\.(Write|Flush|Buffer|BinaryWrite|ContentType|Charset|AddHeader|redirect)
  - score: 2
    value: (?i)(page\.display\(execute\(|eval|Session\(|Public\s+Sub|Private\s+Sub|Public\s+Function)
  - score: -5
    value: \?>$
  - score: -5
    value: $<\?
- language: aspx
  top:
  - score: 100
    value: '(?i)^<%@\s*.*page\s*language="C#".*'
  - score: 100
    value: '(?i)^<%@\s+import\s+Namespace="?System.*"?\s*%>'
  neartop:
  - score: 50
    value: '(?i)^<%@\s+import\s+Namespace="?System.*"?\s*%>'
  - score: 30
    value: '(?i)^<%@\s*.*page.*language="C#".*'	
  - score: 30
    value: '(?i)^<%@\s*.*page.*language="C#".*'
  - score: 10
    value: (?i)codebehind=".*\.cs".*%>
  rules:
  - score: 10
    value: (?i)<asp:Content.*>
  - score: 5
    value: System\.(IO|Data)	
- language: perl
  top:
  - score: 100
    value: '^#!/usr/bin/perl'
  - score: 100
    value: '^#!/usr/bin/env perl'
- language: python
  top:
  - score: 100
    value: '^#!/usr/bin/env python'
  - score: 100
    value: '^#!/usr/bin/python'
  neartop:
  - score: 30
    value: '^#!.*/env\s*python'
  - score: 30
    value: '^#!.*/python'
  rules:
  - score: 5
    value: if\s*__name__\s*==\s*.__main__.
  - score: 5
    value: 'def\s*\w+\s*(.*):'
  - score: 1
    value: ^import\s*\w+$`
