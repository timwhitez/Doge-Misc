using System;
using System.Management;

namespace DisableAdapter
{
 // Token: 0x02000002 RID: 2
 internal class Program
 {
  // Token: 0x06000001 RID: 1 RVA: 0x00002050 File Offset: 0x00000250
  private static void Main(string[] args)
  {
   foreach (ManagementBaseObject managementBaseObject in new ManagementObjectSearcher(new SelectQuery("SELECT * FROM Win32_NetworkAdapter WHERE NetConnectionId != NULL")).Get())
   {
    ManagementObject managementObject = (ManagementObject)managementBaseObject;
    if ((string)managementObject["NetConnectionId"] != "Local Network Connection")
    {
     managementObject.InvokeMethod("Disable", null);
    }
   }
  }
 }
}
