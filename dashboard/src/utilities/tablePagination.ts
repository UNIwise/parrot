export const generatePageNumbers = (pageCount: number, currentPage: number) => {
  const pageNumbers = [];

  // If there are less than 7 pages, show all pages buttons
  // otherwise show only 5 pages buttons
  if (pageCount <= 7) {
    pageNumbers.push(...[...Array(pageCount).keys()].map(i => i++));
  } else {
    if (currentPage <= 4) {
      pageNumbers.push(...[...Array(5).keys()].map(i => i++));
      pageNumbers.push("...", pageCount);
    } else if (currentPage >= pageCount - 3) {
      pageNumbers.push(1, "...");
      pageNumbers.push(...[...Array(5).keys()].map(i => pageCount - 4 + i));
    } else {
      pageNumbers.push(
        1,
        "...",
        currentPage - 1,
        currentPage,
        currentPage + 1,
        "...",
        pageCount,
      );
    }
  }
  return pageNumbers;
};
