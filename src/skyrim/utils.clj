(ns skyrim.utils
  (:require [clojure.string :as str]
            [clojure.set :as set]
            [clojure.data.csv :as csv]
            [clojure.java.io :as io]
            [clojure.math.combinatorics :refer [combinations]]))

(defn sanatize-keyword [input]
  (keyword (str/replace (str/replace (str/lower-case input) #" " "-") #"'" "")))

(def ie-raw
  (with-open [reader (io/reader "alchemy-ingredients.csv")]
    (let [thing #(into #{} (map sanatize-keyword (take 4 (rest %))))
          thing2 (into {} (map (fn [x] {(sanatize-keyword (first x)) (thing x)}) (doall (csv/read-csv reader :separator \;)))) ]
      {:disc thing2 :all thing2})))

(defn eval-potion
  "Evaluate potions determineing the discovered and knowable effects found in each ingredient"
  [[ing1 ing2 ing3] ie-map]
  (let [ing-set-op (fn [x y] (set/intersection (get-in ie-map [:all x]) (get-in ie-raw [:all y])))
        gamma (set/union (ing-set-op ing1 ing2)
                         (ing-set-op ing1 ing3)
                         (ing-set-op ing2 ing3))
        ing-set-op2 (fn [x] (set/intersection (get-in ie-map [:disc x]) gamma))]
    {ing1 (ing-set-op2 ing1)
     ing2 (ing-set-op2 ing2)
     ing3 (ing-set-op2 ing3)}))

(def no-duds
  (let [combos (combinations (keys (:all ie-raw)) 3)]
    (filter #(not= #{} (apply set/union (vals (eval-potion % ie-raw)))) combos)))

(defn effects-left [ie]
  (into {} (filter #(not= #{} (val %)) (:disc ie))))
