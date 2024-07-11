import DarkModeRoundedIcon from "@mui/icons-material/DarkModeRounded";
import LightModeIcon from "@mui/icons-material/LightMode";
import IconButton, { IconButtonProps } from "@mui/joy/IconButton";
import { styled, useColorScheme } from "@mui/joy/styles";
import { useEffect, useState } from "react";

enum ColorMode {
  LIGHT = "light",
  DARK = "dark",
}

const StyledColorSchemeToggle = styled(IconButton)({
  '&[data-mode="dark"]': {
    '& > *:first-of-type': { display: 'none' },
    '& > *:last-child': { display: 'initial' },
  },
  '&[data-mode="light"]': {
    '& > *:first-of-type': { display: 'initial' },
    '& > *:last-child': { display: 'none' },
  },
});

export const ColorSchemeToggle = (props: IconButtonProps) => {
  const { onClick, ...other } = props;
  const { mode, setMode } = useColorScheme();
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
  }, []);

  if (!mounted) {
    return (
      <StyledColorSchemeToggle
        size="sm"
        variant="outlined"
        color="neutral"
        {...other}
        disabled
      />
    );
  }

  return (
    <StyledColorSchemeToggle
      id="toggle-mode"
      size="sm"
      variant="soft"
      color="neutral"
      {...other}
      onClick={(event) => {
        if (mode === ColorMode.LIGHT) {
          setMode(ColorMode.DARK);
        } else {
          setMode(ColorMode.LIGHT);
        }
        onClick?.(event);
      }}
      data-mode={mode}
    >
      <DarkModeRoundedIcon />
      <LightModeIcon />
    </StyledColorSchemeToggle>
  );
};
