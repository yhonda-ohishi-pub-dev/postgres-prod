@echo off
REM Execute harness for Windows
REM Runs claude /execute command without MINGW path conversion issues

setlocal

cd /d "%~dp0.."

:loop
echo.
echo ============================================================
echo Checking plan.md for pending tasks...
echo ============================================================

REM Check if there are pending phases (simple check)
findstr /C:"Phase 1-3:" /C:"Phase 1-4:" /C:"Phase 1-5:" /C:"Phase 1-6:" /C:"Phase 1-7:" plan.md | findstr /V "✅" >nul
if errorlevel 1 (
    echo.
    echo All tasks completed!
    goto :end
)

echo Pending tasks found. Running claude...
echo.

REM Run claude with the execute command
claude --print "/execute"

if errorlevel 1 (
    echo.
    echo Claude execution may have failed.
    set /p continue="Continue? (y/n): "
    if /i not "%continue%"=="y" goto :end
)

REM Small delay before next iteration
timeout /t 2 /nobreak >nul

goto :loop

:end
echo.
echo Execute harness finished.
endlocal
