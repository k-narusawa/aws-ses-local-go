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
      aria-label="ページネーション"
    >
      {page > 1 && (
        <button
          type="button"
          onClick={() => setPage(page - 1)}
          className="inline-flex items-center px-3 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition-colors duration-200"
          aria-label="前のページ"
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
          前へ
        </button>
      )}

      <div className="hidden md:flex space-x-1">
        {Array.from({ length: totalPage }, (_, i) => {
          const pageNum = i + 1;
          const isCurrentPage = pageNum === page;
          const isNearCurrentPage = Math.abs(pageNum - page) <= 2;
          const isFirstPage = pageNum === 1;
          const isLastPage = pageNum === totalPage;

          if (isNearCurrentPage || isFirstPage || isLastPage) {
            return (
              <button
                key={i}
                onClick={() => setPage(pageNum)}
                type="button"
                className={`inline-flex items-center px-4 py-2 text-sm font-medium rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition-colors duration-200 ${
                  isCurrentPage
                    ? "bg-blue-600 text-white"
                    : "text-gray-700 bg-white border border-gray-300 hover:bg-gray-50"
                }`}
                aria-current={isCurrentPage ? "page" : undefined}
                aria-label={`${pageNum}ページ目`}
              >
                {pageNum}
              </button>
            );
          } else if (
            (pageNum === page - 3 && pageNum > 2) ||
            (pageNum === page + 3 && pageNum < totalPage - 1)
          ) {
            return (
              <span
                key={i}
                className="inline-flex items-center px-4 py-2 text-sm text-gray-700"
                aria-hidden="true"
              >
                ...
              </span>
            );
          }
          return null;
        })}
      </div>

      <div className="md:hidden">
        <span className="text-sm text-gray-700">
          {page} / {totalPage} ページ
        </span>
      </div>

      {page < totalPage && (
        <button
          type="button"
          onClick={() => setPage(page + 1)}
          className="inline-flex items-center px-3 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition-colors duration-200"
          aria-label="次のページ"
        >
          次へ
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
