const PLATFORM = {
  MODE: "web",
  NAME: "web",
  ID: "web",
};

type Device =
  | "*"
  | "web"
  | "native"
  | "native:desktop"
  | "native:desktop:windows"
  | "native:desktop:macos"
  | "native:desktop:linux"
  | "native:mobile"
  | "native:mobile:ios"
  | "native:mobile:android";

function processPlatformNameToID(platform: string): Device {
  if (platform === "*") {
    return "*";
  }

  if (platform === "web") {
    return "web";
  }

  if (platform === "native") {
    return "native";
  }

  // Check native platforms
  switch (platform) {
    case "windows":
      return "native:desktop:windows";
    case "macos":
      return "native:desktop:macos";
    case "linux":
      return "native:desktop:linux";
    case "ios":
      return "native:mobile:ios";
    case "android":
      return "native:mobile:android";
    default:
      return "native";
  }
}

export { PLATFORM, processPlatformNameToID, type Device };
