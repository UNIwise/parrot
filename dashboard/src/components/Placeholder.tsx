import { Box, Skeleton } from "@mui/joy";
import { FC } from "react";

export const Placeholder: FC = () => {
  return (
    <Box
      sx={{
        display: "flex",
        flexDirection: "column",
        justifyContent: "space-around",
      }}
    >
      <Skeleton variant="rectangular" height={70} />
      <Box
        sx={{
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
          height: 150,
          marginTop: 25,
          marginBottom: 25,
        }}
      >
        <Skeleton sx={{ borderRadius: "50%", width: 150, height: 150, position: "relative" }} />
      </Box>
      <Skeleton variant="rectangular" height={70} />
    </Box>
  );
};
