import IconButton, { IconButtonProps } from "@mui/joy/IconButton";
import { useColorScheme } from "@mui/joy/styles";
import { useEffect, useState } from "react";

import DarkModeRoundedIcon from "@mui/icons-material/DarkModeRounded";
import LightModeIcon from "@mui/icons-material/LightMode";

enum colorMode {
  LIGHT = "light",
  DARK = "dark",
}

export const ColorSchemeToggle = (props: IconButtonProps) => {
  const { onClick, sx, ...other } = props;
  const { mode, setMode } = useColorScheme();
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
  }, []);

  if (!mounted) {
    return (
      <IconButton
        size="sm"
        variant="outlined"
        color="neutral"
        {...other}
        sx={sx}
        disabled
      />
    );
  }

  return (
    <IconButton
      id="toggle-mode"
      size="sm"
      variant="soft"
      color="neutral"
      {...other}
      onClick={(event) => {
        if (mode === colorMode.LIGHT) {
          setMode(colorMode.DARK);
        } else {
          setMode(colorMode.LIGHT);
        }
        onClick?.(event);
      }}
      sx={[
        mode === colorMode.DARK
          ? { "& > *:first-of-type": { display: "none" } }
          : { "& > *:first-of-type": { display: "initial" } },
        mode === colorMode.LIGHT
          ? { "& > *:last-child": { display: "none" } }
          : { "& > *:last-child": { display: "initial" } },
        ...(Array.isArray(sx) ? sx : [sx]),
      ]}
    >
      <DarkModeRoundedIcon />
      <LightModeIcon />
    </IconButton>
  );
};
