import KeyboardArrowLeftIcon from "@mui/icons-material/KeyboardArrowLeft";
import KeyboardArrowRightIcon from "@mui/icons-material/KeyboardArrowRight";
import { Box, Button, IconButton, iconButtonClasses } from "@mui/joy";
import { FC } from "react";
import { generatePageNumbers } from "../utilities/tablePagination";

type TablePaginationSectionProps = {
  currentPage: number;
  pageCount: number;
  onPageChange: (page: number) => void;
}

export const TablePaginationSection: FC<TablePaginationSectionProps> = ({
  currentPage,
  pageCount,
  onPageChange,
}) => {
  return (
    <Box
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
        color="primary"
        startDecorator={<KeyboardArrowLeftIcon />}
        onClick={() => onPageChange(Math.max(1, currentPage - 1))}
        disabled={currentPage === 1}
      >
        Previous
      </Button>

      <Box sx={{ flex: 1 }} />
      {generatePageNumbers(pageCount, currentPage).map((page, index) => (
        <IconButton
          key={index}
          size="sm"
          variant={typeof page === 'number' ? "outlined" : "plain"}
          color="primary"
          onClick={() => typeof page === 'number' && onPageChange(page)}
          disabled={page === currentPage || typeof page !== 'number'}
        >
          {page}
        </IconButton>
      ))}
      <Box sx={{ flex: 1 }} />

      <Button
        size="sm"
        variant="outlined"
        color="primary"
        endDecorator={<KeyboardArrowRightIcon />}
        onClick={() => onPageChange(Math.min(pageCount, currentPage + 1))}
        disabled={currentPage === pageCount}
      >
        Next
      </Button>
    </Box>
  );
};
