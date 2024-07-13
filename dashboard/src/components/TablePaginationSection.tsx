import KeyboardArrowLeftIcon from "@mui/icons-material/KeyboardArrowLeft";
import KeyboardArrowRightIcon from "@mui/icons-material/KeyboardArrowRight";
import { Box, Button, IconButton, iconButtonClasses } from "@mui/joy";

interface TablePaginationSectionProps {
  currentPage: number;
  pageCount: number;
  onPageChange: (page: number) => void;
}

export const TablePaginationSection: React.FC<TablePaginationSectionProps> = ({
  currentPage,
  pageCount,
  onPageChange,
}) => {
  const generatePageNumbers = () => {
    const pageNumbers = [];
    if (pageCount <= 7) {
      for (let i = 1; i <= pageCount; i++) {
        pageNumbers.push(i);
      }
    } else {
      if (currentPage <= 4) {
        for (let i = 1; i <= 5; i++) {
          pageNumbers.push(i);
        }
        pageNumbers.push("...", pageCount);
      } else if (currentPage >= pageCount - 3) {
        pageNumbers.push(1, "...");
        for (let i = pageCount - 4; i <= pageCount; i++) {
          pageNumbers.push(i);
        }
      } else {
        pageNumbers.push(1, "...", currentPage - 1, currentPage, currentPage + 1, "...", pageCount);
      }
    }
    return pageNumbers;
  };

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
        color="neutral"
        startDecorator={<KeyboardArrowLeftIcon />}
        onClick={() => onPageChange(Math.max(1, currentPage - 1))}
        disabled={currentPage === 1}
      >
        Previous
      </Button>

      <Box sx={{ flex: 1 }} />
      {generatePageNumbers().map((page, index) => (
        <IconButton
          key={index}
          size="sm"
          variant={typeof page === 'number' ? "outlined" : "plain"}
          color="neutral"
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
        color="neutral"
        endDecorator={<KeyboardArrowRightIcon />}
        onClick={() => onPageChange(Math.min(pageCount, currentPage + 1))}
        disabled={currentPage === pageCount}
      >
        Next
      </Button>
    </Box>
  );
};
