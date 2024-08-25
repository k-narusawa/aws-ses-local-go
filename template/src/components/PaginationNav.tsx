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
    <>
      <nav className="flex items-center space-x-1" aria-label="Pagination">
        {page > 1 && (
          <button
            type="button"
            onClick={() => setPage(page - 1)}
            className="p-2.5 min-w-[40px] inline-flex justify-center items-center gap-x-2 text-sm rounded-full text-gray-800 hover:bg-gray-100 focus:outline-none focus:bg-gray-100 disabled:opacity-50 disabled:pointer-events-none"
            aria-label="Previous"
          >
            <span className="sr-only">Previous</span>
            <span aria-hidden="true">«</span>
          </button>
        )}

        {Array.from({ length: totalPage }, (_, i) => (
          <button
            key={i}
            onClick={() => setPage(i + 1)}
            type="button"
            className={`min-w-[40px] flex justify-center items-center text-gray-800 hover:bg-gray-100 focus:outline-none focus:bg-gray-100 py-2.5 text-sm rounded-full disabled:opacity-50 disabled:pointer-events-none ${
              i + 1 === page ? "bg-gray-100" : ""
            }`}
            aria-current="page"
          >
            {i + 1}
          </button>
        ))}

        {page < totalPage && (
          <button
            type="button"
            onClick={() => setPage(page + 1)}
            className="p-2.5 min-w-[40px] inline-flex justify-center items-center gap-x-2 text-sm rounded-full text-gray-800 hover:bg-gray-100 focus:outline-none focus:bg-gray-100 disabled:opacity-50 disabled:pointer-events-none"
            aria-label="Next"
          >
            <span className="sr-only">Next</span>
            <span aria-hidden="true">»</span>
          </button>
        )}
      </nav>
    </>
  );
};
