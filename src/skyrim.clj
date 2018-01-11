(ns skyrim
  (:require [clojure.string :as str]
            [clojure.data.csv :as csv]
            [clojure.java.io :as io]
            [clojure.algo.generic.functor :refer [fmap]]
            [clojure.walk :refer [keywordize-keys]]
            [clojure.math.combinatorics :refer [combinations partitions]]))

(def ingredient-effects-map
  (with-open [reader (io/reader "alchemy-ingredients.csv")]
    (into {}
          (map
           (fn [x] {(first x) (take 4 (rest x))})
           (doall
            (csv/read-csv reader :separator \;))))))

(def remaining (atom ingredient-effects-map))

(defn make-potion
  "takes as argument a list of ingredients and returns a set of the effects"
  [ingredient-seq & {:keys [use-remaining removed-effects remove-remaining] :or {use-remaining nil removed-effects nil remove-remaining nil}}]
  (let [remaining-effects (map set (vals (select-keys @remaining use-remaining))) ;; [set]
        full-effects (map set (vals (select-keys ingredient-effects-map ingredient-seq)))
        effs (concat remaining-effects full-effects)
        potion (apply clojure.set/union
                      (concat (map clojure.set/intersection (map set (combinations effs 2))))) ;;set[set]
        rem-d (map #(clojure.set/intersection % (set (apply concat potion))) remaining-effects)]
    (if remove-remaining
      (doseq [ing use-remaining
             eff potion]
        (swap! remaining (fn [x] (update x ing #(remove #{eff} %)))))
      nil)
    (if removed-effects
      rem-d
      potion)))


(defn f [x]
  (let [[q p] x]
    (make-potion p :use-remaining q :removed-effects true)))
