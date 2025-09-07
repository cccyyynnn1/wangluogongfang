import {
  URLEncode,
  URLDecode,
  EncodeText,
  DecodeText,
  ToBase64,
  FromBase64,
  ToHex,
  FromHex,
} from '@/assets/cyberChef/core/operations/index.mjs'

import _utils from '@/assets/cyberChef/core/Utils.mjs'

enum targetType {
  to = 0,
  from,
}

type formattingType = 'ASCII' | 'UTF-8' | 'GBK' | 'GB2312' | 'UTF-16BE' | 'UTF-16LE' | 'URL' | 'HEX' | 'BASE64'

type excludetype = Exclude<formattingType, 'URL' | 'HEX' | 'BASE64'>

type params = {
  data: any // 源数据
  formatting: formattingType // 数据类型
  target: targetType // to:编码;from:解码
  arg?: string[] // 剩余参数
}

const DecodeOrEncodeTextCode: formattingType[] = ['ASCII', 'UTF-8', 'GBK', 'GB2312', 'UTF-16BE', 'UTF-16LE']

const methods: Record<formattingType, Function[]> = {
  ASCII: [EncodeText, DecodeText],
  'UTF-8': [EncodeText, DecodeText],
  GBK: [EncodeText, DecodeText],
  GB2312: [EncodeText, DecodeText],
  'UTF-16BE': [EncodeText, DecodeText],
  'UTF-16LE': [EncodeText, DecodeText],
  URL: [URLEncode, URLDecode],
  HEX: [ToHex, FromHex],
  BASE64: [ToBase64, FromBase64],
}

const CHR_ENC_CODE_PAGES: Record<excludetype, string> = {
  ASCII: 'US-ASCII (7-bit) (20127)',
  'UTF-8': 'UTF-8 (65001)',
  GBK: 'MAC Simplified Chinese (GB 2312) (10008)',
  GB2312: 'Simplified Chinese GB2312 (20936)',
  'UTF-16BE': 'UTF-16BE (1201)',
  'UTF-16LE': 'UTF-16LE (1200)',
}
/**
 * @description '编码解码'
 * @params  { data: any ,formatting: formattingType ,target: 'to' | 'from' ,arg?: string[] }
 * @return string
 */
const data1 = `HTTP/1.1 200
Server: nginx
Content-Type: text/html; charset=GB2312
Content-Length: 7578
Connection: keep-alive
Set-Cookie: JSESSIONID=0000aYw90Tc7SN3o432R1y3YZ8x:-1; Path=/; HttpOnly
Date:Thu,08Aug202405:49:14GMTX-Powered-By:Servlet/3.1Content-Language:en-USExpires:Thu,01Dec199416:00:00GMTCache-Control:no-cache="set-cookie,set-cookie2"
<html><METAhttp-equiv="Content-Style-Type"content="text/css"><LINKhref="/Issue_sys/theme/Master.css"rel="stylesheet"type="text/css"><scripttype="text/javascript">window.onload=function(){varform=document.getElementById("form1");form.appId.value=ccdcJsFunction.getBrowserInfo("lv1_menu_id");form.clientId.value=ccdcJsFunction.getBrowserInfo("client_id");form.userId.value=ccdcJsFunction.getBrowserInfo("user_id");form.token.value=ccdcJsFunction.getBrowserInfo("token");form.channelType.value=ccdcJsFunction.getBrowserInfo("login_channel");form.loginName.value=ccdcJsFunction.getBrowserInfo("user_login_name");form.userZHName.value=encodeURI(JSON.parse(ccdcJsFunction.getBrowserInfo("user_info"))["user_name"]);form.clientName.value=encodeURI(JSON.parse(ccdcJsFunction.getBrowserInfo("client_info"))["client_name"]);var_bindAcctList=eval(ccdcJsFunction.getBrowserInfo("bind_acct_list"));varbindAcctList=[];for(vari=0;i<_bindAcctList.length;i++){if(_bindAcctList[i].acct_type==="01"){bindAcctList.push(_bindAcctList[i]);}}//���ʻ���������ת��varmap={};map['01']="ծȯ";map['02']="����";map['03']="����";map['04']="�ʽ�";map['05']="��ֵ";map['06']="��ҵծ";//�洢�������˵��˺���ϢvarshowAcctList=[];//ѡ��һ��ʱ���ø�����ֵ������setAccount=function(i){if(showAcctList.length>0){form.bondAcct.value=showAcctList[i]["bond_acct"];form.bondAcctName.value=encodeURI(showAcctList[i]["grp_acct_no_name"]);}}//����һ�������ݸ�ʽ�����б�setAccountList=function(accts){vartable=document.getElementById("multiAccounts");for(vari=0;i<accts.length;i++){vartr=table.insertRow();vartd=tr.insertCell()td.setAttribute("class","table-text")td.innerHTML='<inputTYPE="radio"NAME="selection"onclick="setAccount('+i+')"</input>';vartd2=tr.insertCell();td2.setAttribute("class","table-text");td2.innerHTML=accts[i]['bond_acct'];vartd3=tr.insertCell();td3.setAttribute("class","table-text");td3.innerHTML=accts[i]['grp_acct_no_name'];vartd4=tr.insertCell();td4.setAttribute("class","table-text");//td4.innerHTML=accts[i]['biz_line'];td4.innerHTML=map[accts[i]['biz_line']]}}checkAccount=function(data){if(data.length==0){alert("û�дӵ�½��Ϣ�л�ȡ����Ӫ�˺�");return;}varreqdata='params='+JSON.stringify(data);varxhr=newXMLHttpRequest();xhr.open('POST','/Issue_sys/check.do',false);xhr.setRequestHeader('Content-Type','application/x-www-form-urlencoded');xhr.onreadystatechange=function(){if(xhr.readyState==4){if(xhr.status==200||xhr.status==304){varresp=JSON.parse(xhr.responseText);if(resp!=null&&resp.hasOwnProperty("length")&&resp.length==0){alert("δ��ȡ����Ч�˻���");}//ʹ�÷��ص��б����˲�ѯ����//showAcctList=bindAcctList.filter(item=>resp.indexOf(item["grp_acct_no"])!=-1);showAcctList=resp;setAccountList(resp);//����Ĭ��ѡ���б���һ��varradios=document.getElementsByName("selection");if(radios.length>0){radios[0].click();}}else{alert("��ѯ�˻��ʸ���������,������:"+xhr.status);}}};xhr.send(reqdata);}checkAccount(ccdcJsFunction.getBrowserInfo("user_id"));//checkAccount(bindAcctList.map(item=>item['grp_acct_no']));//�ύ������ť��������ͬ��summitForm=function(){form.submit();}}</SCRIPT><formid="form1"action="init.do"method="post"><inputname="appId"type="hidden"/><inputname="clientId"type="hidden"/><inputname="userId"type="hidden"/><inputname="token"type="hidden"/><inputname="channelType"type="hidden"/><inputname="loginName"type="hidden"/><inputname="userZHName"type="hidden"/><inputname="bondAcct"type="hidden"/><inputname="bondAcctName"type="hidden"/><inputname="clientName"type="hidden"/></form><TABLEalign=centerborder=0cellPadding=0cellSpacing=0width=760><TBODY><TR><TD><TABLEborder=0cellPadding=0cellSpacing=0width=760><TBODY><TR><TD><IMGalt="���������˷���"height=82src="/Issue_sys/img/logo_qt.jpg"width=760></TD></TR></TBODY></TABLE></TD></TR></TBODY></TABLE><TABLEalign="center"class="noframe-framing-table"border="0"cellpadding="3"cellspacing="0"width="760"><TBODY><TRvalign="top"><TDclass="top-navigation">��ҳ&nbsp;&nbsp;|&nbsp;&nbsp;����ծȯ����ϵͳ&nbsp;&nbsp;|&nbsp;&nbsp;�����г�����ϵͳ&nbsp;&nbsp;|&nbsp;&nbsp;ע��&nbsp;&nbsp;|<fontonclick="window.alert('���������˷���ϵͳV1.0016(20100903)')">&nbsp;&nbsp;����</font>&nbsp;&nbsp;|&nbsp;&nbsp;</TD><TDclass="top-navigation"align="right">&nbsp;&nbsp;</TD></TR></TBODY></TABLE><BR><TABLEALIGN="CENTER"BORDER="0"CELLSPACING="0"CELLPADDING="10"WIDTH="760"><TBODY><TR><TDCLASS="layout-manager"ID="notabs"><inputtype="hidden"name="command"value="Issue_main"/><inputtype="hidden"name="selected_index"value="-1"/><inputtype="hidden"name="ZBSXH"value=""/><tableborder="0"cellpadding="5"cellspacing="0"valign="top"width="100%"><trvalign="top"><tdclass="function-button-section"nowrap><inputtype="button"name="button_enter"value="����"class="buttons"id="functions"onclick="summitForm()"/></td></tr></table><TABLEid="multiAccounts"align="center"border="0"cellspacing="1"cellpadding="0"width="100%"class="framing-table"><TBODY><TR><TH>ѡ��</TH><TH>�˻�</TH><TH>�˻���</TH><TH>ҵ������</TH></TR></TBODY></TABLE></TD></TR></TBODY></TABLE></html>`
const data2 = `HTTP/1.1 200
Server: Apache/2.4.39 (Unix)
Content-Type: text/html; charset=gbk
Content-Length: 9070
Connection: Keep-Alive
Date:Wed,22May201901:04:18GMTX-Powered-By:PHP/7.3.4Expires:Thu,19Nov198108:52:00GMTCache-Control:no-store,no-cache,must-revalidatePragma:no-cacheKeep-Alive:timeout=5,max=99
<html><head><metahttp-equiv="Content-Type"content="text/html;charset=gbk"><title>editfile-172.29.5.9</title><styletype="text/css">body,td{font:12pxArial,Tahoma;line-height:16px;}.input,select{font:12pxArial,Tahoma;background:#fff;border:1pxsolid#666;padding:2px;height:22px;}.area{font:12px'CourierNew',Monospace;background:#fff;border:1pxsolid#666;padding:2px;}.red{color:#f00;}.black{color:#000;}.green{color:#090;}.b{font-weight:bold;}.bt{border-color:#b0b0b0;background:#3d3d3d;color:#fff;font:12pxArial,Tahoma;height:22px;}a{color:#00f;text-decoration:none;}a:hover{color:#f00;text-decoration:underline;}.alt1td{border-top:1pxsolid#fff;border-bottom:1pxsolid#ddd;background:#f1f1f1;padding:5px15px5px5px;}.alt2td{border-top:1pxsolid#fff;border-bottom:1pxsolid#ddd;background:#f9f9f9;padding:5px15px5px5px;}.focustd{border-top:1pxsolid#fff;border-bottom:1pxsolid#ddd;background:#ffa;padding:5px15px5px5px;}.headtd{border-top:1pxsolid#fff;border-bottom:1pxsolid#ddd;background:#e9e9e9;padding:5px15px5px5px;font-weight:bold;}.headtdspan{font-weight:normal;}.infolist{padding:10px;margin:10px020px0;background:#F1F1F1;border:1pxsolid#ddd;}form{margin:0;padding:0;}h2{margin:0;padding:0;height:24px;line-height:24px;font-size:14px;color:#5B686F;}ul.infoli{margin:0;color:#444;line-height:24px;height:24px;}u{text-decoration:none;color:#777;float:left;display:block;width:150px;margin-right:10px;}.drives{padding:5px;}.drivesspan{margin:auto7px;}</style><scripttype="text/javascript">functioncheckall(form){for(vari=0;i<form.elements.length;i++){vare=form.elements[i];if(e.type=='checkbox'){if(e.name!='chkall'&&e.name!='saveasfile')e.checked=form.chkall.checked;}}}function$(id){returndocument.getElementById(id);}functioncreatedir(){varnewdirname;newdirname=prompt('Pleaseinputthedirectoryname:','');if(!newdirname)return;g(null,null,'createdir',newdirname);}functionfileperm(pfile,val){varnewperm;newperm=prompt('Currentdir/file:'+pfile+'\nPleaseinputnewpermissions:',val);if(!newperm)return;g(null,null,'fileperm',pfile,newperm);}functionrename(oldname){varnewfilename;newfilename=prompt('Filename:'+oldname+'\nPleaseinputnewfilename:','');if(!newfilename)return;g(null,null,'rename',newfilename,oldname);}functioncreatefile(){varfilename;filename=prompt('Pleaseinputthefilename:','');if(!filename)return;g('editfile',null,null,filename);}functionsetdb(dbname){if(!dbname)return;$('dbform').tablename.value='';$('dbform').doing.value='';if($('dbform').sql_query){$('dbform').sql_query.value='';}$('dbform').submit();}functionsetsort(k){$('dbform').order.value=k;$('dbform').submit();}functionsettable(tablename,doing){if(!tablename)return;if(doing){$('dbform').doing.value=doing;}else{$('dbform').doing.value='';}$('dbform').sql_query.value='';$('dbform').tablename.value=tablename;$('dbform').submit();}functions(act,cwd,p1,p2,p3,p4,charset){if(act!=null)$('opform').act.value=act;if(cwd!=null)$('opform').cwd.value=cwd;if(p1!=null)$('opform').p1.value=p1;if(p2!=null)$('opform').p2.value=p2;if(p3!=null)$('opform').p3.value=p3;if(p4!=null){$('opform').p4.value=p4;}else{$('opform').p4.value='';}if(charset!=null)$('opform').charset.value=charset;}functiong(act,cwd,p1,p2,p3,p4,charset){s(act,cwd,p1,p2,p3,p4,charset);$('opform').submit();}</script></head><bodystyle="margin:0;table-layout:fixed;word-break:break-all"><formname="opform"id="opform"action="/shell.php"method="post"><inputid="act"type="hidden"name="act"value="editfile"/><inputid="cwd"type="hidden"name="cwd"value="/html/web/"/><inputid="p1"type="hidden"name="p1"value=""/><inputid="p2"type="hidden"name="p2"value="index.html"/><inputid="p3"type="hidden"name="p3"value=""/><inputid="p4"type="hidden"name="p4"value=""/><inputid="charset"type="hidden"name="charset"value="gbk"/></form><tablewidth="100%"border="0"cellpadding="0"cellspacing="0"><trclass="head"><td><spanstyle="float:right;">Linuxnt4.15.0-47-generic#50-UbuntuSMPWedMar1310:44:52UTC2019x86_64/User:0()/Group:0()</span>172.29.5.9(172.29.5.9)</td></tr><trclass="alt1"><td><spanstyle="float:right;">Charset:<selectclass="input"id="charset"name="charset"onchange="g(null,null,null,null,null,null,this.value);"><optionvalue="big5">big5</option><optionvalue="cp-866">cp866</option><optionvalue="euc-jp">ujis</option><optionvalue="euc-kr">euckr</option><optionvalue="gbk"selected>gbk</option><optionvalue="iso-8859-1">latin1</option><optionvalue="koi8-r">koi8r</option><optionvalue="koi8-u">koi8u</option><optionvalue="utf-8">utf8</option><optionvalue="windows-1252">latin1</option></select></span><ahref="javascript:g('logout');">Logout</a>|<ahref="javascript:g('file',null,'','','','','gbk');">FileManager</a>|<ahref="javascript:g('mysqladmin',null,'','','','','gbk');">MYSQLManager</a>|<ahref="javascript:g('shell',null,'','','','','gbk');">ExecuteCommand</a>|<ahref="javascript:g('phpenv',null,'','','','','gbk');">PHPVariable</a>|<ahref="javascript:g('portscan',null,'','','','','gbk');">PortScan</a>|<ahref="javascript:g('secinfo',null,'','','','','gbk');">Securityinformation</a>|<ahref="javascript:g('eval',null,'','','','','gbk');">EvalPHPCode</a>|<ahref="javascript:g('backconnect',null,'','','','','gbk');">BackConnect</a></td></tr></table><tablewidth="100%"border="0"cellpadding="15"cellspacing="0"><tr><td><formname="form1"id="form1"action="/shell.php"method="post"onsubmit="g('editfile',null,'edit',this.p2.value,this.p3.value);returnfalse;"><h2>Create/EditFile&raquo;</h2><p>Filename<br/><inputclass="input"name="p2"id="p2"value="index.html"type="text"size="100"/></p><p>FileContent<br/><textareaclass="area"id="p3"name="p3"cols="100"rows="25">﻿&lt;!DOCTYPEhtmlPUBLIC&quot;-//W3C//DTDXHTML1.0Transitional//EN&quot;&quot;http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd&quot;&gt;&lt;htmlxmlns=&quot;http://www.w3.org/1999/xhtml&quot;&gt;&lt;!----&gt;&lt;head&gt;&lt;metahttp-equiv=&quot;Content-Type&quot;content=&quot;text/html;charset=utf-8&quot;/&gt;&lt;title&gt;&lt;/title&gt;&lt;linkhref=&quot;static/cwym.css&quot;rel=&quot;stylesheet&quot;type=&quot;text/css&quot;/&gt;&lt;/head&gt;&lt;body&gt;&lt;divclass=&quot;cwym&quot;&gt;&lt;divclass=&quot;cwym_nr&quot;&gt;&lt;divclass=&quot;cwym_top&quot;&gt;&lt;divclass=&quot;logo&quot;&gt;&lt;imgsrc=&quot;static/logo.gif&quot;/&gt;&lt;/div&gt;&lt;/div&gt;&lt;divclass=&quot;cwym_middle&quot;&gt;&lt;divclass=&quot;cwym_middlenr&quot;&gt;&lt;div&gt;亲爱的朋友：&lt;br/&gt;您可以从以下查询入口查询战队信息&lt;br/&gt;&lt;span&gt;您所需要做的是：篡改index.html页面的内容信息（文字信息和logo信息）。&lt;/span&gt;&lt;/div&gt;&lt;div&gt;Dearfriends:&lt;br/&gt;Youcanquerytheteaminformationfromthefollowingentrylink.&lt;br/&gt;Allyouneedtodoistotamperwithcontentinformation(textandlogo)oftheindex.htmlpage.&lt;br/&gt;&lt;/div&gt;&lt;/div&gt;&lt;divclass=&quot;button&quot;style=&quot;font-size:15px;&quot;&gt;&lt;ahref=&quot;show.php&quot;style=&quot;background-image:url('static/btn.gif');height:200px;width:100px;&quot;&gt;query&lt;/a&gt;&lt;/div&gt;&lt;/div&gt;&lt;/div&gt;&lt;/body&gt;&lt;/html&gt;</textarea></p><p><inputclass="bt"name="submit"id="submit"type="submit"value="Submit"></p></form><formaction="/shell.php"method="post"><inputtype="hidden"name="act"value="file"/><inputtype="hidden"name="cwd"value="/html/web/"/><inputtype="hidden"name="charset"value="gbk"/><p><inputclass="bt"type="submit"value="Goback..."></p></form></td></tr></table><divstyle="padding:10px;border-bottom:1pxsolid#fff;border-top:1pxsolid#ddd;background:#eee;"><spanstyle="float:right;">Processedin0.000353second(s)</span>Poweredby<atitle="Build20190516"href="http://www.mimicSoftwareTeam.net"target="_blank">MimicSoftwareTeam2019final</a>.Copyright(C)2019<ahref="http://www.mimicSoftwareTeam.net"target="_blank">[S4T]</a>AllRightsReserved.</div></body></html>`
const useCodec = (payload: params): string => {
  const { data, formatting, target, arg } = payload
  const method = methods[formatting][target]
  let args: any[] = []
  // const data1 = '你好'
  if (DecodeOrEncodeTextCode.includes(formatting)) {
    const code = CHR_ENC_CODE_PAGES[formatting as excludetype]
    args = arg?.length ? [code, ...arg!] : [code]
  }
  // @ts-ignore
  const op = new method()
  return op.run(data, args)
}

export default useCodec
