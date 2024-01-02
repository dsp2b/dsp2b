# 一些逆向记录

## 一些蓝图示例

### 只有一个传送带
BLUEPRINT:0,10,0,0,0,0,0,0,638395809227381327,0.10.28.21014,%E6%96%B0%E8%93%9D%E5%9B%BE,"H4sIAAAAAAAAC2NkYGBgZEAAJiQ+I8N/BoYTUGFGsBQYdPyrtwfLH5DchsxuvrbZCYQvsisz/IcCZOPhjFXq880YJ0iaI7N7NLc4gTBIMwiANDMi6QEALJ3gGqkAAAA="909071918E2F234FAB35E92CD054F866

### 只有一个熔炉
BLUEPRINT:0,10,0,0,0,0,0,0,638395809227381327,0.10.28.21014,%E6%96%B0%E8%93%9D%E5%9B%BE,"H4sIAAAAAAAAC2NkQAWMUAxh/2dgOAFlMsKFEWoPSG7Dxv7HYcfwHwpQTWZgAAB4dngncAAAAA=="2881F7A76BAF3A19C17C948A5C773D72

BLUEPRINT:0,10,0,0,0,0,0,0,638395809227381327,0.10.28.21014,新蓝图,"H4sIAAAAAAAAC2NkQAWMUAxh/2dgOAFlMsKFEWoPSG7Dxv7HYcfwHwpQTWZgAAB4dngncAAAAA=="2881F7A76BAF3A19C17C948A5C773D72


# BlueprintData

蓝图头

```c#
	public string headerStr
	{
		get
		{
			string[] array = new string[21];
			array[0] = "BLUEPRINT:0,";// 固定标识
			int num = 1;
			int num2 = (int)this.layout;
			array[num] = num2.ToString();// 10
			array[2] = ",";
			array[3] = this.icon0.ToString();// 图标1
			array[4] = ",";
			array[5] = this.icon1.ToString();// 图标2
			array[6] = ",";
			array[7] = this.icon2.ToString();// 图标3
			array[8] = ",";
			array[9] = this.icon3.ToString();// 图标4
			array[10] = ",";
			array[11] = this.icon4.ToString();// 图标5
			array[12] = ",0,";
			array[13] = this.time.Ticks.ToString();// 时间 (638395809227381327-621355968000000000)/10000/1000=unix时间戳
			array[14] = ",";
			array[15] = this.gameVersion;// 游戏版本
			array[16] = ",";
			array[17] = this.shortDesc.Escape();// 蓝图文件url编码
			array[18] = ",";
			array[19] = this.desc.Escape();// 蓝图描述url编码
			array[20] = "\"";
			return string.Concat(array);
		}
	}
```

新蓝图

```c#
	// Token: 0x0600122D RID: 4653 RVA: 0x00138000 File Offset: 0x00136200
	public void ResetAsEmpty()
	{
		this.time = DateTime.Now;
		this.gameVersion = GameConfig.gameVersion.ToFullString();
		this.shortDesc = "新的蓝图".Translate();
		this.desc = "";
		this.layout = EIconLayout.OneIcon;
		this.icon0 = 0;
		this.icon1 = 0;
		this.icon2 = 0;
		this.icon3 = 0;
		this.icon4 = 0;
		this.cursorOffset_x = 0;
		this.cursorOffset_y = 0;
		this.dragBoxSize_x = 1;
		this.dragBoxSize_y = 1;
		this.cursorTargetArea = 0;
		this.primaryAreaIdx = 0;
		this.areas = new BlueprintArea[1];
		this.areas[0] = new BlueprintArea();
		this.areas[0].index = 0;
		this.areas[0].parentIndex = -1;
		this.areas[0].tropicAnchor = 0;
		this.areas[0].areaSegments = 200;
		this.areas[0].anchorLocalOffsetX = 0;
		this.areas[0].anchorLocalOffsetY = 0;
		this.areas[0].width = 1;
		this.areas[0].height = 1;
		this.buildings = new BlueprintBuilding[0];
	}
```

蓝图编码输出

```c#
	// Token: 0x06001239 RID: 4665 RVA: 0x00138C2C File Offset: 0x00136E2C
	public string ToBase64String()
	{
		try
		{
			StringBuilder stringBuilder = new StringBuilder(1024);
			using (MemoryStream memoryStream = new MemoryStream())
			{
				using (BinaryWriter binaryWriter = new BinaryWriter(memoryStream))
				{
					this.Export(binaryWriter);
					memoryStream.Position = 0L;
					using (MemoryStream memoryStream2 = new MemoryStream())
					{
						using (GZipStream gzipStream = new GZipStream(memoryStream2, CompressionMode.Compress))
						{
							memoryStream.CopyTo(gzipStream);
						}
						byte[] inArray = memoryStream2.ToArray();
						stringBuilder.Append(this.headerStr);
						stringBuilder.Append(Convert.ToBase64String(inArray));
						string value = MD5F.Compute(stringBuilder.ToString());
						stringBuilder.Append("\"");
						stringBuilder.Append(value);
						return stringBuilder.ToString();
					}
				}
			}
		}
		catch (Exception ex)
		{
			Debug.LogError(ex.ToString().Replace("Exception", "Excption"));
		}
		return "";
	}
```

导出建筑数据
```c#
	public void Export(BinaryWriter w)
	{
		this.version = 1;
		w.Write(this.version);
		w.Write(this.cursorOffset_x);
		w.Write(this.cursorOffset_y);
		w.Write(this.cursorTargetArea);
		w.Write(this.dragBoxSize_x);
		w.Write(this.dragBoxSize_y);
		w.Write(this.primaryAreaIdx);
		int num = (this.areas != null) ? this.areas.Length : 0;
		w.Write((byte)num);
		for (int i = 0; i < num; i++)
		{
			this.areas[i].Export(w);
		}
		int num2 = (this.buildings != null) ? this.buildings.Length : 0;
		w.Write(num2);
		for (int j = 0; j < num2; j++)
		{
			this.buildings[j].Export(w);
		}
	}
```

## 资源导出

使用AssetStudio和AssetBundleExtractor工具在\DSPGAME_Data\sharedassets0.assets中导出

## 物品、公式

数据存储在：LDB中ItemProtoSet为物品，RecipeProtoSet为公式，资源文件为：Prototypes/xxx

ItemProto.Preload