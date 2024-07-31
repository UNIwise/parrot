import { Sheet, Typography } from "@mui/joy";
import { FC } from "react";
import { Link } from "react-router-dom";
import { ColorSchemeToggle } from "./ColorSchemeToggle";

export const Header: FC = () => {
  return (
    <Sheet
      component="header"
      sx={{
        p: 2,
        background: `linear-gradient(
                  to right,
                  rgba(255, 0, 0, 0.7),
                  rgba(255, 165, 0, 0.7),
                  rgba(255, 255, 0, 0.7),
                  rgba(0, 128, 0, 0.7),
                  rgba(0, 0, 255, 0.7),
                  rgba(128, 0, 128, 0.7)
                )`,
        mb: 6,
        display: "flex",
        justifyContent: "center",
        borderRadius: "5px",
        alignItems: "center",
        height: 64,
      }}
    >
      <Typography
        component={Link}
        to={"/"}
        level="h4"
        sx={{
          fontWeight: "bold",
          color: "primary.solidColor",
          fontSize: "2.2em",
        }}
      >
        Parrot
      </Typography>
      <ColorSchemeToggle sx={{ ml: "auto" }} />
    </Sheet>
  );
};
