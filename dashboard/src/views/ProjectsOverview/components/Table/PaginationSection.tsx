import KeyboardArrowLeftIcon from "@mui/icons-material/KeyboardArrowLeft";
import KeyboardArrowRightIcon from "@mui/icons-material/KeyboardArrowRight";
import { Box, Button, IconButton, iconButtonClasses } from "@mui/joy";

export const PaginationSection = () => {
  //TODO: Implement pagination logic. Try to use cached data from query hook using query key
  return (
    <Box
      className="Pagination-laptopUp"
      sx={{
        pt: 2,
        gap: 1,
        [`& .${iconButtonClasses.root}`]: { borderRadius: "50%" },
        display: {
          xs: "none",
          md: "flex",
        },
      }}
    >
      <Button
        size="sm"
        variant="outlined"
        color="neutral"
        startDecorator={<KeyboardArrowLeftIcon />}
      >
        Previous
      </Button>

      <Box sx={{ flex: 1 }} />
      {["1", "2", "3", "…", "8", "9", "10"].map((page) => (
        <IconButton
          key={page}
          size="sm"
          variant={Number(page) ? "outlined" : "plain"}
          color="neutral"
        >
          {page}
        </IconButton>
      ))}
      <Box sx={{ flex: 1 }} />
      <Button
        size="sm"
        variant="outlined"
        color="neutral"
        endDecorator={<KeyboardArrowRightIcon />}
      >
        Next
      </Button>
    </Box>
  );
};
