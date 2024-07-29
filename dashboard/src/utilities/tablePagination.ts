export const generatePageNumbers = (pageCount: number, currentPage: number, ) => {
  const pageNumbers = [];

  // If there are less than 7 pages, show all pages buttons
  // otherwise show only 5 pages buttons
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
