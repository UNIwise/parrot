import { Box, CircularProgress } from "@mui/joy";
import { FC } from "react";

export const Placeholder: FC = () => {
  return (
    <Box
      sx={{
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        height: 150,
        marginTop: "15%",
      }}
    >
      <CircularProgress size="sm" value={35}>
        <CircularProgress size="md" value={60}>
          <CircularProgress size="lg" value={70} />
        </CircularProgress>
      </CircularProgress>
    </Box>
  );
};
