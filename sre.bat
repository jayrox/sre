@ECHO off
REM Example usage
REM Loops from 1 to 8000 searching Sonarr for missing episodes by index.
REM Change 8000 to whatever MAX you'd like, ideally close to the total missing episodes in your library.

:STARTSCRIPT
    FOR /L %%G IN (1,1,8000) DO (
        cls
        ECHO %%G
        sre -i %%G
        ECHO.
        TIMEOUT /T 20
    )

GOTO :STARTSCRIPT
