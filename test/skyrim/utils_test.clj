(ns skyrim.utils-test
  (:require [clojure.test :refer [deftest testing is] :as t]
            [skyrim.utils :refer [eval-potion no-duds ie-raw]]))

(deftest test-eval-potion
  (testing "Testing if the most valueable potion is valid"
    (is (= {:wheat #{:damage-stamina-regen :fortify-health}
            :creep-cluster #{:damage-stamina-regen :fortify-carry-weight}
            :giants-toe #{:damage-stamina-regen :fortify-carry-weight :fortify-health}}
           (eval-potion [:giants-toe :wheat :creep-cluster] ie-raw)) )))

(deftest test-no-duds
  (testing "That non-duds is working"
    (is (not= 0 (count no-duds)))))
