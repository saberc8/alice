export enum BasicStatus {
	DISABLE = 0,
	ENABLE = 1,
}

export enum ResultStatus {
	SUCCESS = 0,
	ERROR = -1,
	TIMEOUT = 401,
}

export enum StorageEnum {
	UserInfo = "userInfo",
	UserToken = "userToken",
	Settings = "settings",
}

export enum ThemeMode {
	Light = "light",
	Dark = "dark",
}

export enum ThemeLayout {
	Vertical = "vertical",
	Horizontal = "horizontal",
	Mini = "mini",
}

export enum ThemeColorPresets {
	Default = "default",
	Cyan = "cyan",
	Purple = "purple",
	Blue = "blue",
	Orange = "orange",
	Red = "red",
}

export enum LocalEnum {
	zh_CN = "zh_CN",
}

export enum PermissionType {
	GROUP = 0,
	CATALOGUE = 1,
	MENU = 2,
	COMPONENT = 3,
}

export enum HtmlDataAttribute {
	ColorPalette = "data-color-palette",
	ThemeMode = "data-theme-mode",
}
