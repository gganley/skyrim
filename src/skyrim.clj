(ns skyrim
  (:require [clojure.data.csv :as csv]
            [clojure.java.io :as io]))

(def ie-map
  (with-open [reader (io/reader "alchemy-ingredients.csv")]
    (into {}
          (map
           (fn [x] {(first x) (set (take 4 (rest x)))})
           (doall
            (csv/read-csv reader :separator \;))))))
