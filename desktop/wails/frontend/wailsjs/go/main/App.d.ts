// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {main} from '../models';

export function AddBlock(arg1:string,arg2:string,arg3:string,arg4:string,arg5:number):Promise<Record<string, any>>;

export function AddBlockAtPath(arg1:string,arg2:main.BlockPath,arg3:string):Promise<Record<string, any>>;

export function CaptureDOM(arg1:string,arg2:string):Promise<void>;

export function CaptureNavigationHistory(arg1:Array<string>):Promise<void>;

export function GetAsset(arg1:string):Promise<string>;

export function GetBacklinks(arg1:string):Promise<Record<string, Array<string>>>;

export function GetPage(arg1:string):Promise<main.PageData>;

export function GetPageList():Promise<Array<string>>;

export function InitTestCapture():Promise<void>;

export function IsTestMode():Promise<boolean>;

export function LoadDirectory(arg1:string):Promise<void>;

export function LogResourceError(arg1:Record<string, any>):Promise<void>;

export function LogUserAction(arg1:string):Promise<void>;

export function RefreshPages():Promise<void>;

export function UpdateBlock(arg1:string,arg2:string,arg3:string):Promise<void>;

export function UpdateBlockAtPath(arg1:string,arg2:main.BlockPath,arg3:string):Promise<Record<string, any>>;
