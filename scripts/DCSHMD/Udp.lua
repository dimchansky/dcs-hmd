local socket = require("socket")

DCSHMD_Udp = {}

DCSHMD_Udp.Host = "127.0.0.1"
DCSHMD_Udp.Port = 19089

DCSHMD_Udp.Socket = nil

-- Simulation id
DCSHMD_Udp.SimID = string.format("%08x*", os.time())

-- State data for export
DCSHMD_Udp.PacketSize = 0
DCSHMD_Udp.SendStrings = {}
DCSHMD_Udp.LastData = {}

-- Frame counter for non important data
DCSHMD_Udp.TickCount = 0

function DCSHMD_Udp.Start()
    if DCSHMD_Udp.Socket ~= nil then
        DCSHMD_Udp.Socket:close()
        DCSHMD_Udp.Socket = nil
    end

    DCSHMD_Udp.PacketSize = 0
    DCSHMD_Udp.SendStrings = {}
    DCSHMD_Udp.LastData = {}

    DCSHMD_Udp.Socket = socket.udp()

    if DCSHMD_Udp.Socket == nil then return false end

    DCSHMD_Udp.Socket:setsockname("*", 0)
    DCSHMD_Udp.Socket:setoption('broadcast', true)
    DCSHMD_Udp.Socket:settimeout(.001) -- set the timeout for reading the socket

    return true
end

function DCSHMD_Udp.Stop()
    if DCSHMD_Udp.Socket ~= nil then
        DCSHMD_Udp.Socket:close()
    end
end

function DCSHMD_Udp.Receive()
    if DCSHMD_Udp.Socket == nil then return nil end

    return DCSHMD_Udp.Socket:receive()
end

function DCSHMD_Udp.Send(id, value)
    if string.len(value) > 3 and value == string.sub("-0.00000000", 1, string.len(value)) then
        value = value:sub(2)
    end

    if DCSHMD_Udp.LastData[id] == nil or DCSHMD_Udp.LastData[id] ~= value then
        local data =  id .. "=" .. value
        local dataLen = string.len(data)

        if dataLen + DCSHMD_Udp.PacketSize > 576 then
            DCSHMD_Udp.Flush()
        end

        table.insert(DCSHMD_Udp.SendStrings, data)

        DCSHMD_Udp.LastData[id] = value
        DCSHMD_Udp.PacketSize = DCSHMD_Udp.PacketSize + dataLen + 1
    end
end

function DCSHMD_Udp.Flush()
    if #DCSHMD_Udp.SendStrings > 0 then
        local packet = DCSHMD_Udp.SimID..table.concat(DCSHMD_Udp.SendStrings, ":").."\n"

        socket.try(DCSHMD_Udp.Socket:sendto(packet, DCSHMD_Udp.Host, DCSHMD_Udp.Port))

        DCSHMD_Udp.SendStrings = {}
        DCSHMD_Udp.PacketSize = 0
    end
end

function DCSHMD_Udp.ResetChangeValues()
    DCSHMD_Udp.LastData = {}
    DCSHMD_Udp.TickCount = 10
end
