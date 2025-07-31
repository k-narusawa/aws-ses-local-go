import React from "react";

type Props = {
  page: number;
  totalPage: number;
  setPage: (page: number) => void;
};

export const PaginationNav: React.FC<Props> = ({
  page,
  totalPage,
  setPage,
}) => {
  return (
    <nav
      className="flex items-center justify-center space-x-2"
      aria-label="Pagination"
    >
      {page > 1 && (
        <button
          type="button"
          onClick={() => setPage(page - 1)}
          className="inline-flex items-center px-3 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition-colors duration-200"
          aria-label="Previous page"
        >
          <svg
            className="w-5 h-5 mr-1"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M15 19l-7-7 7-7"
            />
          </svg>
          Previous
        </button>
      )}

      <div className="hidden md:flex space-x-1">
        {(() => {
          const visiblePages = new Set<number>();

          // Always show first and last page
          visiblePages.add(1);
          visiblePages.add(totalPage);

          // Add current page and 4 pages before and after
          for (
            let i = Math.max(2, page - 4);
            i <= Math.min(totalPage - 1, page + 4);
            i++
          ) {
            visiblePages.add(i);
          }

          // Convert to sorted array
          const sortedPages = Array.from(visiblePages).sort((a, b) => a - b);

          return sortedPages.map((pageNum, index) => {
            const isCurrentPage = pageNum === page;

            // Show ellipsis if there's a gap
            if (index > 0 && sortedPages[index - 1] !== pageNum - 1) {
              return (
                <React.Fragment key={`${pageNum}-ellipsis`}>
                  <span
                    className="inline-flex items-center px-4 py-2 text-sm text-gray-700"
                    aria-hidden="true"
                  >
                    ...
                  </span>
                  <button
                    onClick={() => setPage(pageNum)}
                    type="button"
                    className={`inline-flex items-center px-4 py-2 text-sm font-medium rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition-colors duration-200 ${
                      isCurrentPage
                        ? "bg-blue-600 text-white"
                        : "text-gray-700 bg-white border border-gray-300 hover:bg-gray-50"
                    }`}
                    aria-current={isCurrentPage ? "page" : undefined}
                    aria-label={`Page ${pageNum}`}
                  >
                    {pageNum}
                  </button>
                </React.Fragment>
              );
            }

            return (
              <button
                key={pageNum}
                onClick={() => setPage(pageNum)}
                type="button"
                className={`inline-flex items-center px-4 py-2 text-sm font-medium rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition-colors duration-200 ${
                  isCurrentPage
                    ? "bg-blue-600 text-white"
                    : "text-gray-700 bg-white border border-gray-300 hover:bg-gray-50"
                }`}
                aria-current={isCurrentPage ? "page" : undefined}
                aria-label={`Page ${pageNum}`}
              >
                {pageNum}
              </button>
            );
          });
        })()}
      </div>

      <div className="md:hidden">
        <span className="text-sm text-gray-700">
          Page {page} of {totalPage}
        </span>
      </div>

      {page < totalPage && (
        <button
          type="button"
          onClick={() => setPage(page + 1)}
          className="inline-flex items-center px-3 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition-colors duration-200"
          aria-label="Next page"
        >
          Next
          <svg
            className="w-5 h-5 ml-1"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M9 5l7 7-7 7"
            />
          </svg>
        </button>
      )}
    </nav>
  );
};
