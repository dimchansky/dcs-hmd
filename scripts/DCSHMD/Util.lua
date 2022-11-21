DCSHMD_Util = {}

function DCSHMD_Util.CompareString(str1, str2)
    if str1 == str2 or string.find(str1, str2) ~= nil then
        return true
    end

    return false
end
