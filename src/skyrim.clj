(ns skyrim
  (:require [clojure.core.async :as async]
            [clojure.java.io :as io]
            [clojure.data.csv :as csv]))


(defn sanatize-keyword [^String input]
  (keyword (clojure.string/replace (clojure.string/lower-case input) #" " "-")))

(def ie-raw
  (with-open [reader (io/reader "alchemy-ingredients.csv")]
    (map (fn [x] {(sanatize-keyword (first x)) (into #{} (map sanatize-keyword (take 4 (rest x))))}) (doall (csv/read-csv reader :separator \;)))))
