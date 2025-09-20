# Premier League Winner Predictor

![Python](https://img.shields.io/badge/python-3.9%2B-blue)
![Streamlit](https://img.shields.io/badge/streamlit-app-red)
![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)
![Maintainer](https://img.shields.io/badge/maintainer-deethecreator-blue)


**Initiator:** `deethecreator`

This project downloads historical **Premier League match results** from **football-data.co.uk**, engineers per-team, per-season features, and trains **four ML models** to predict the **league champion**:

- Logistic Regression (baseline)
- Random Forest
- Gradient Boosting (XGBoost)
- Neural Network (MLP)

It also includes a **Streamlit dashboard** to inspect feature importance and run quick what-if predictions.

---

## 🚀 Quickstart

```bash
git clone <your-repo>
cd premier-league-predictor
python -m venv .venv && source .venv/bin/activate  # on Windows: .venv\Scripts\activate
pip install -r requirements.txt

# 1) Download raw CSVs (2010-11 to 2023-24, and latest season file if available)
python src/download_data.py

# 2) Build features per team-season
python src/preprocess.py

# 3) Train & compare 4 models (GroupKFold by season), pick best and save to disk
python src/train.py

# 4) (Optional) Evaluate saved model on a holdout or partial-season snapshot
python src/evaluate.py

# 5) Run dashboard
streamlit run dashboard/app.py
```

> Note: Steps 1–3 will create `data/processed/epl_features.csv` and save the best model to `data/processed/best_model.pkl` with a JSON report in `data/processed/model_report.json`.

---

## 📁 Structure

```
premier-league-predictor/
├── data/
│   ├── raw/          # downloaded CSVs
│   └── processed/    # engineered features, models, reports
├── src/
│   ├── download_data.py
│   ├── preprocess.py
│   ├── train.py
│   └── evaluate.py
├── dashboard/
│   └── app.py
├── notebooks/
│   └── model_training.ipynb
├── requirements.txt
└── README.md
```

---

## 🧠 Features (per team-season)

- Matches Played, Wins, Draws, Losses
- Goals For, Goals Against, Goal Difference
- Points (3/win, 1/draw), Points/Game
- Shots, Shots on Target (if available)
- Corners (if available)
- First-Half/Second-Half goal splits (if available)
- Simple recent-form proxy (last-5 match points average) aggregated for the season

**Label:** `Champion` (1 for the season winner, else 0), using points → GD → GF tiebreakers.

---

## 📊 Model Comparison
We compare 4 models via **GroupKFold** (grouped by `Season`) to avoid leakage across seasons. Metrics: **Accuracy**, **F1 (macro)**, **ROC-AUC**.

The training script saves:
- `data/processed/best_model.pkl` (scikit-learn pipeline)
- `data/processed/model_report.json` (metrics per model and overall ranking)

---

## ⚠️ Data Source
- Historical CSVs: [football-data.co.uk](https://www.football-data.co.uk/englandm.php)
- Filenames like `E0_2010_2011.csv`, `E0_2011_2012.csv`, …, `E0_2023_2024.csv` and current season `E0.csv` (if present).

---

## 🧪 Repro Tips
- If some columns are missing in older seasons (e.g., `xG`), the pipeline auto-fills or drops gracefully.
- Renaming of teams across years is lightly normalized (e.g., "Man United" vs "Manchester United").


## Documentation
- [Docs Overview](docs/README.md)
- [Architecture](docs/architecture.md)
- [Data](docs/data.md)
- [Modeling](docs/modeling.md)
- [Dashboard](docs/dashboard.md)
- [How to Use](docs/how_to_use.md)


## 🔎 Preview

![Dashboard Preview](docs/images/dashboard.png)
