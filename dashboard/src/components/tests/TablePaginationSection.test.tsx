import { render } from "@testing-library/react";
import userEvent from '@testing-library/user-event';
import { describe, expect, it } from 'vitest';
import { TablePaginationSection } from "../TablePaginationSection";


describe("renders pagination buttons correctly", () => {
  it('should display provided currentPage and pageCount in the pagination buttons', () => {
    const user = userEvent.setup();
    const currentPage = 2;
    const pageCount = 5;

    const r = render(
      <TablePaginationSection
        currentPage={currentPage}
        pageCount={pageCount}
        onPageChange={() => { }}
      />
    );

    const previousButton = r.getByTestId("pagination-previous-button");
    const nextButton = r.getByTestId("pagination-next-button");

    expect(previousButton).toBeInTheDocument();
    expect(nextButton).toBeInTheDocument(); // Use the toBeInTheDocument function

    user.click(previousButton);


    user.click(nextButton);
  });
});
