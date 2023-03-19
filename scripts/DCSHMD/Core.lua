dofile(lfs.writedir()..[[Scripts\DCSHMD\Util.lua]])
dofile(lfs.writedir()..[[Scripts\DCSHMD\Udp.lua]])

DCSHMD = {}

DCSHMD.DebugFile = nil
DCSHMD.Interval = 0.01 -- frequency of export events (sec)

function DCSHMD.Start()
    DCSHMD.DebugFile = io.open(lfs.writedir()..[[Logs\DCSHMDDebug.log]], "wa")

    log.write('DCSHMD EXPORT',log.INFO,'Mission Started')

    if not DCSHMD_Udp.Start() then
        if DCSHMD.DebugFile then
            DCSHMD.DebugFile:write("ERROR CREATING DCSHMD SOCKET!\r\n")
        end

        return
    end

    if DCSHMD.DebugFile then
        DCSHMD.DebugFile:write("DCSHMD SOCKET CREATED!\r\n")
    end
end

function DCSHMD.BeforeNextFrame()
end

function DCSHMD.ActivityNextEvent()
    DCSHMD_Udp.TickCount = DCSHMD_Udp.TickCount + 1

    local selfdata = LoGetSelfData()

    -- Check if we are on an aircraft
    if selfdata == nil then return end

    if DCSHMD_Util.CompareString(selfdata.Name, "Ka-50") then
        local lDevice = GetDevice(0)

        if type(lDevice) == "table" then

            lDevice:update_arguments()

            -- Handle the simple-case data that can be simply read via device:get_argument_value
            DCSHMD.ProcessArguments(lDevice, DCSHMD.Ka50HighImportanceArguments)

            DCSHMD_Udp.Flush()
        end
    end
end

function DCSHMD.Stop()
    DCSHMD_Udp.Stop()

    log.write('DCSHMD EXPORT',log.INFO,'Mission Ended')

    if DCSHMD.DebugFile then
        DCSHMD.DebugFile:write("LuaExportStop called.\r\n")
        io.close(DCSHMD.DebugFile)
    end
end

-- Handles simple-case data that can be simply read via device:get_argument_value
function DCSHMD.ProcessArguments(device, arguments)
    if arguments == nil then return end

    local lArgument , lFormat , lArgumentValue

    -- Other airplanes
    for lArgument, lFormat in pairs(arguments) do
        lArgumentValue = string.format(lFormat,device:get_argument_value(lArgument))
       DCSHMD_Udp.Send(lArgument, lArgumentValue)
    end
    --end
end

DCSHMD.Ka50HighImportanceArguments =
{
    -- VVI
    ---------------------------------------------------
    [24]  = "%.4f", 		-- vy (Vertical Velocity Indicator) input={-30.0, 30.0} output={-1.0,1.0}
    -- Rotor Pitch
    ---------------------------------------------------
    [53]  = "%.4f", 		-- RotorPitch input={1.0, 15.0} output={0.0,1.0}
    -- Rotor RPM
    ---------------------------------------------------
    [52]  = "%.4f"   		-- RotorRPM input={0.0, 110.0} output={0.0,1.0}
}