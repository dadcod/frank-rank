import { GridStack, GridStackOptions, GridStackWidget } from "gridstack";

const granularityMultiplier = 2;
let ranksGrid: GridStack;
let rankerGrid: GridStack;
let poolGrid: GridStack;
let id = 0;

export function setupGridstack() {
  ranksGrid = initRanksGrid();
  rankerGrid = initRankerGrid();
  poolGrid = initPoolGrid();

  addGridEvents();
  addEventHandlers();
  setupMutationObserver();
  ranksGrid.cellHeight(rankerGrid.getCellHeight());
}

function initRanksGrid(): GridStack {
  const rows: GridStackWidget[] = [
    {
      x: 0,
      y: 0,
      w: granularityMultiplier,
      h: granularityMultiplier,
      content: "row 1",
    },
    {
      x: 0,
      y: granularityMultiplier,
      w: granularityMultiplier,
      h: granularityMultiplier,
      content: "row 2",
    },
  ];

  let ranksGrid = GridStack.init(
    {
      resizable: { handles: "s, n" },
      margin: 0.5,
      float: true,
      removable: true,
      column: granularityMultiplier,
    },
    ".ranks"
  );
  ranksGrid.load(rows);
  return ranksGrid;
}
function initRankerGrid() {
  return GridStack.init(
    {
      resizable: { handles: "s, n" },
      margin: 5,
      float: true,
      acceptWidgets: true,
      removable: true,
      column: 24,
      minRow: granularityMultiplier * 2,
    },
    ".ranker"
  );
}

function initPoolGrid() {
  var pool: GridStackWidget[] = [
    {
      x: 0,
      y: 0,
      w: granularityMultiplier,
      h: granularityMultiplier,
      content: "pool 1",
      id: `${id++}`,
    },
  ];
  var poolGrid = GridStack.init(
    {
      resizable: { handles: "s, n" },
      margin: 5,
      float: true,
      removable: true,
      minRow: granularityMultiplier,
      column: 24,
    },
    ".pool"
  );
  poolGrid.load(pool);
  return poolGrid;
}

function addEventHandlers() {
  document.getElementById("add-item")?.addEventListener("click", () => {
    poolGrid.compact();
    const row =
      Math.trunc(
        (poolGrid.getGridItems().length * granularityMultiplier) /
          poolGrid.getColumn()
      ) * granularityMultiplier;
    console.log(row);
    const col =
      (poolGrid.getGridItems().length * granularityMultiplier) %
      poolGrid.getColumn();
    console.log(col);
    const node: GridStackWidget = {
      x: col,
      y: row,
      w: granularityMultiplier,
      h: granularityMultiplier,
      id: `${id++}`,
    };
    poolGrid.addWidget(
      `<div><div class="grid-stack-item-content"><div class="custom-content">New Item ${id}</div></div></div>`,
      node
    );
  });

  document.getElementById("add-row")?.addEventListener("click", () => {
    var node = {
      x: 0,
      y: ranksGrid
        .getGridItems()
        .reduce((acc, row) => acc + (row.gridstackNode?.h || 0), 0),
      w: granularityMultiplier,
      h: granularityMultiplier,
    };
    ranksGrid.addWidget(
      '<div><div class="grid-stack-item-content">New Item</div></div>',
      node
    );
  });

  document
    .getElementById("change-columns")
    ?.addEventListener("click", function () {
      const items = rankerGrid.save(false);
      console.log(items);
      rankerGrid.column(rankerGrid.getColumn() + 1, "none");
      // setTimeout(() => {
      //   granularityMultiplier += 1;
      //   rankerGrid.load([...items, ...rows]);
      //   addRow();
      // }, 500);
    });
  document.getElementById("save")?.addEventListener("click", function () {
    console.log(rankerGrid.save(true, true));
  });
  document.getElementById("sort")?.addEventListener("click", function () {
    rankerGrid.compact("compact", true);
  });
}

function setupMutationObserver() {
  // Use MutationObserver to detect changes in the grid container
  var observerConfig = {
    attributes: true,
    childList: true,
    subtree: true,
  };

  var observer = new MutationObserver((_, __) => {
    setGridLines(
      ranksGrid.el,
      rankerGrid.cellWidth(),
      rankerGrid.getCellHeight()
    );
    setGridLines(poolGrid.el, poolGrid.cellWidth(), poolGrid.getCellHeight());
    setGridLines(rankerGrid.el, rankerGrid.cellWidth(), rankerGrid.cellWidth());
    ranksGrid.cellHeight(rankerGrid.getCellHeight());
  });

  observer.observe(poolGrid.el, observerConfig);
  observer.observe(rankerGrid.el, observerConfig);
  observer.observe(ranksGrid.el, observerConfig);
}
// Function to set the background grid lines dynamically
function setGridLines(
  gridElement: HTMLElement,
  cellWidth: number,
  cellHeight: number
) {
  gridElement.style.backgroundSize = `${cellWidth}px ${cellHeight}px`;
  gridElement.style.backgroundImage = `linear-gradient(to right, lightgray 1px, transparent 1px),
                                              linear-gradient(to bottom, lightgray 1px, transparent 1px)`;
}

function adjustGridRows() {
  const savedVersion = rankerGrid.save(true, true);
  rankerGrid.destroy(false);
  console.log(ranksGrid.getRow());
  rankerGrid = GridStack.init(
    {
      ...savedVersion,
      minRow: ranksGrid.getRow(),
    } as GridStackOptions,
    "ranker"
  );
}

function addGridEvents() {
  ranksGrid.on("change", adjustGridRows);
  ranksGrid.on("added", adjustGridRows);

  rankerGrid.on("removed", function (_, items) {
    const removed = `<div>${items[0].el?.innerHTML}</div>`;
    poolGrid.addWidget(removed, {
      x: 0,
      y: 0,
      w: granularityMultiplier,
      h: granularityMultiplier,
    });
    poolGrid.compact();
  });

  rankerGrid.on("added", function () {
    setGridLines(
      rankerGrid.el,
      rankerGrid.cellWidth(),
      rankerGrid.getCellHeight()
    );
  });

  poolGrid.on("added", function () {
    setGridLines(poolGrid.el, poolGrid.cellWidth(), poolGrid.getCellHeight());
  });
}
