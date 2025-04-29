revision_id: 4128175

文件路径
```text
offboard/perception/tracking/python/scripts/auto_create_and_upload_scenario.py
```

文件内容
```python
·1: #!/usr/bin/env python3
·2: """Auto create and upload scenario.
·3: 
·4: Required Arguments:
·5: --trip segment id: The id of trip segment that we want to create scenario.
·6: --start-timestamp.
·7: --end-timestamp
·8: --experiment-id
·9: --effective-timestamp
·10: --issue-type.
·11: 
·12: Sample usage:
·13: rosrun offboard_perception_tracking auto_create_and_upload_scenario.py \
·14:  --issue-id cn12333917 \
·15:     --trip-id 17161_20240614_095537 \
·16:         --start-timestamp 1718330335027 \
·17:           --end-timestamp 1718330345027 \
·18:           --experiment-id offboard.detection.tracking \
·19:             --effective-timestamp 1718330342027 --issue-type Velocity
·20: """
·21: import argparse
·22: import json
·23: import os
·24: import subprocess
·25: import tempfile
+26: from voy_tempest import bag
·27: 
·28: from google.protobuf import json_format, text_format
·29: 
·30: from scenario_tools import auto_create_scenario_utils
·31: from sim_common.db_accessor import ScenarioAccessor
·32: from sim_protos import scenario_pb2
·33: from voy_base.cli_formatter import CliFormatter
·34: from voy_data_utils.regions import Regions
·35: from voy_protos import module_type_pb2 as voy
·36: 
+37: _DOWNLOAD_BAG_PATH = "/tmp/auto_create_trip_segment.bag"
```
>行号：37-37
>
>评论：硬编码的路径`/tmp/auto_create_trip_segment.bag`可能会导致文件冲突或权限问题。
>
>建议：使用`tempfile`模块生成临时文件路径，以避免硬编码路径带来的问题。
>```diff
>- _DOWNLOAD_BAG_PATH = "/tmp/auto_create_trip_segment.bag"
>+ _DOWNLOAD_BAG_PATH = tempfile.mktemp(prefix='auto_create_trip_segment', suffix='.bag')
>```
>人工标注:[Helpful or Unhelpful]
```python
+38: 
·39: 
·40: def _get_args():
·41:   """Defines the command line arguments and returns the parsed arguments."""
·42:   parser = argparse.ArgumentParser(
·43:       description='Uploads the auto-created scenario to Trail')
·44:   parser.add_argument(
·45:       '--issue-id', type=str, help='Issue id', required=True)
·46:   parser.add_argument(
·47:       '--trip-id', type=str, help='Tripsegment id', required=True)
·48:   parser.add_argument(
·49:       '--experiment-id', type=str, help='Experiment id', required=True)
·50:   parser.add_argument(
·51:       '--start-timestamp',
·52:       type=int,
·53:       help='The start timestamp of trip segment',
·54:       required=True)
·55:   parser.add_argument(
·56:       '--end-timestamp',
·57:       type=int,
·58:       help='The end timestamp of trip segment',
·59:       required=True)
·60:   parser.add_argument(
·61:       '--effective-timestamp',
·62:       type=int,
·63:       help='The end timestamp of trip segment',
·64:       required=True)
·65:   parser.add_argument(
·66:       '--issue-type',
·67:       type=str,
·68:       help='The issue type of this trip segment',
·69:       required=True)
·70:   args = parser.parse_args()
·71:   return args
·72: 
·73: 
·74: def _auto_create_scenario_and_save(args, scenario_output_name):
·75:   """Auto create scenario and save it to |scenario_output_name| file."""
·76:   run_auto_create_track_scenario = CliFormatter(
·77:       command=
·78:       'rosrun offboard_perception_tracking \
·79:       auto_create_track_sequential_metric_scenario'
·80:   )
·81:   run_auto_create_track_scenario.add('--issue_id', args.issue_id)
·82:   run_auto_create_track_scenario.add('--trip_id', args.trip_id)
·83:   run_auto_create_track_scenario.add('--start_timestamp', args.start_timestamp)
·84:   run_auto_create_track_scenario.add('--end_timestamp', args.end_timestamp)
·85:   run_auto_create_track_scenario.add('--output_file', scenario_output_name)
·86:   run_auto_create_track_scenario.add('--experiment_id',
·87:                                      "offboard.detection.tracking")
- ·88:   run_auto_create_track_scenario.add('--effective_timestamp',
+ ·89:   run_auto_create_track_scenario.add('--effective_timestamp',
·89:                                      args.effective_timestamp)
·90:   run_auto_create_track_scenario.add('--issue_type', args.issue_type)
·91: 
·92:   run_auto_create_track_scenario.add(
·93:       '--s3_config_file', "perception_ofs_config.txt")
·94:   run_auto_create_track_scenario.add('--s3_path_prefix', "offline_sensing/prt")
+95:   run_auto_create_track_scenario.add('--bag_path', _DOWNLOAD_BAG_PATH)
·96:   subprocess.check_output(str(run_auto_create_track_scenario), shell=True)
·97: 
·98: 
·99: def _upload_scenario(args, filename):
·100:   """Uploads scenario result to Trail."""
·101:   assert os.path.isfile(filename)
·102:   scenario = scenario_pb2.Scenario()
·103:   with open(filename) as f:
·104:     text_format.Parse(f.read(), scenario)
·105: 
·106:   scenario.name = auto_create_scenario_utils.AUTO_CREATE_SCENARIO_NAME.format(
·107:       trip_id=args.trip_id,
·108:       issue_time=args.effective_timestamp,
·109:       scenario_topic_name=args.issue_type)
·110: 
·111:   scenario.warmup_ms = 3000
·112: 
·113:   scenario.enabled_modules.extend([
·114:       voy.ModuleType.SENSING, voy.ModuleType.PERCEPTION_TRACKING,
·115:       voy.ModuleType.PERFECT_POSE
·116:   ])
·117: 
·118:   scenario_dict = json_format.MessageToDict(scenario)
·119: 
·120:   # MessageToDict() will convert int64 to string, thus, we need to convert
·121:   # string to int. If we do not convert it, metric frame will be wrong.
·122:   for metric in scenario_dict['metrics']:
·123:     if 'trackSequentialMetric' in metric:
·124:       for metric_frame in metric['trackSequentialMetric']['metricFrames']:
·125:         metric_frame['startTimestamp'] = int(metric_frame['startTimestamp'])
·126: 
·127:   # Define a list of labels for the scenario.
·128:   # The labels are created based on the issue type and issue id.
·129:   labels = [
·130:       "{}.auto_create_scenario".format(args.issue_type), "#{}".format(
·131:           args.issue_id[2:])
·132:   ]
·133: 
·134:   # TODO(torettomarui): add scenario label or tags for this scenario.
·135:   result = ScenarioAccessor.add(
·136:       dict(
·137:           scenario=json.dumps(scenario_dict),
·138:           name=scenario.name,
·139:           labels=",".join(labels),
·140:           updater="torettomarui"),
·141:       region=Regions.CN)
·142:   print(result['id'])
·143: 
·144: 
·145: def main():
·146:   """Main function."""
·147:   args = _get_args()
+148: 
+149:   reader = bag.BagReader(issue_id=f'{args.issue_id}')
+150:   reader.save_to_bag(f'{_DOWNLOAD_BAG_PATH}')
```
>行号：150-150
>
>评论：在`main`函数中，`reader.save_to_bag`方法的调用没有处理可能的异常，可能会导致程序崩溃。
>
>建议：使用`try-except`块来捕获并处理可能的异常，以提高代码的健壮性。
>```diff
>- reader.save_to_bag(f'{_DOWNLOAD_BAG_PATH}')
>+ try:
>+     reader.save_to_bag(f'{_DOWNLOAD_BAG_PATH}')
>+ except Exception as e:
>+     print(f'Error saving bag: {e}')
>+     return
>```
>人工标注:[Helpful or Unhelpful]
```python
+151: 
·152:   scenario_output = tempfile.NamedTemporaryFile(
·153:       prefix='auto_scenario', delete=True)
·154:   _auto_create_scenario_and_save(args, scenario_output.name)
·155:   _upload_scenario(args, scenario_output.name)
·156: 
·157: 
·158: if __name__ == '__main__':
·159:   main()
```

