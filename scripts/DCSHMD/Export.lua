package.path = package.path..";.\\LuaSocket\\?.lua;"
package.cpath = package.cpath..";.\\LuaSocket\\?.dll;"

dofile(lfs.writedir()..[[Scripts\DCSHMD\Core.lua]])

local PrevExport = {}
PrevExport.LuaExportStart = LuaExportStart
PrevExport.LuaExportStop = LuaExportStop
PrevExport.LuaExportActivityNextEvent = LuaExportActivityNextEvent
PrevExport.LuaExportBeforeNextFrame = LuaExportBeforeNextFrame

log.write('DCSHMD EXPORT',log.INFO,'DCSHMD export initialized')

function LuaExportStart()
    local status, err = pcall(function()
        DCSHMD.Start()
    end)

    if not status then
        log.write('ERROR DCSHMD.Start', log.INFO, err)
    end

    if PrevExport.LuaExportStart then
        PrevExport.LuaExportStart()
    end
end

function LuaExportStop()
    local status, err = pcall(function()
        DCSHMD.Stop()
    end)

    if not status then
        log.write('ERROR DCSHMD.Stop', log.INFO, err)
    end

    if PrevExport.LuaExportStop then
        PrevExport.LuaExportStop()
    end
end

function LuaExportBeforeNextFrame()
    local status, err = pcall(function()
        DCSHMD.BeforeNextFrame()
    end)

    if not status then
        log.write('ERROR DCSHMD.BeforeNextFrame', log.INFO, err)
    end

    if PrevExport.LuaExportBeforeNextFrame then
        PrevExport.LuaExportBeforeNextFrame()
    end

    return NextTime
end

function LuaExportActivityNextEvent(currenttime)
    local NextTime = currenttime + DCSHMD.Interval

    local status, err = pcall(function()
        DCSHMD.ActivityNextEvent()
    end)

    if not status then
        log.write('ERROR DCSHMD.ActivityNextEvent', log.INFO, err)
    end

    if PrevExport.LuaExportActivityNextEvent then
        PrevExport.LuaExportActivityNextEvent(currenttime)
    end

    return NextTime
end

